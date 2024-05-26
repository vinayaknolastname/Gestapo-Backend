package service

import (
	"context"
	"fmt"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/service/sso"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (auth *authenticationService) SendOTP(ctx context.Context, req *proto.SendOTPRequest) (*proto.Response, error) {
	err := validateSendOTPRequest(req)
	if err != nil {
		auth.log.LogError("Error while validateSendOTPRequest", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	column, value := helpers.IdentifiesColumnValue(req.GetEmail(), req.GetPhone())
	if req.Action == utils.SIGN_UP {
		if len(column) == 0 {
			auth.log.LogError("Error while IdentifiesColumnValue", column)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		res, err := auth.storage.CheckDataExist(column, value)
		if err != nil {
			auth.log.LogError("Error while CheckDataExist", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		if res {
			err = fmt.Errorf("account already exist using this %s", column)
			auth.log.LogError(err)
			return nil, status.Errorf(codes.AlreadyExists, "account already exist using this %s", column)
		}
	}

	if !helpers.IsEmpty(req.Email) {
		err = auth.emailService.SendOTP(req.Email, utils.EmailSubject, utils.EmailSubject, auth.redis)
		if err != nil {
			auth.log.LogError("Error while SendOTP", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	} else {
		phoneNumber := fmt.Sprintf("+91%s", req.Phone)
		err = auth.twilioService.SendOTP(phoneNumber)
		if err != nil {
			auth.log.LogError("Error while SendOTP", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}
	sessionToken, err := auth.token.CreateSessionToken(value, req.Action)
	if err != nil {
		auth.log.LogError("Error while CreateSessionToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "OTP sent successfully",
	}

	mdOut := metadata.New(map[string]string{
		"session-token": sessionToken,
	})
	return response, grpc.SetHeader(ctx, mdOut)
}

func (auth *authenticationService) verifyOTP(payload *token.SessionPayload, email, phone, code, action string) (bool, error) {
	// auth.log.LogInfo(payload.TokenType)
	// if payload.TokenType != action {
	// 	auth.log.LogError("Payload doesnot match")
	// 	return false, status.Errorf(codes.PermissionDenied, "Unauthorized: Payload doesnot match")
	// }

	column, value := helpers.IdentifiesColumnValue(email, phone)
	if action == utils.SIGN_UP {
		if len(column) == 0 {
			auth.log.LogError("Error while IdentifiesColumnValue", column)
			return false, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		res, err := auth.storage.CheckDataExist(column, value)
		if err != nil {
			auth.log.LogError("Error while CheckDataExist", err)
			return false, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		if res {
			fmtError := fmt.Errorf("account already exist using this %s", column)
			auth.log.LogError(fmtError)
			return false, status.Errorf(codes.AlreadyExists, "account already exist using this %s", column)
		}
	}

	if !helpers.IsEmpty(email) {
		if payload.Value != email {
			auth.log.LogError("Forbidden")
			return false, status.Errorf(codes.PermissionDenied, "Forbidden")
		}
		sts, err := auth.emailService.VerfiyOTP(email, code, auth.redis)
		if err != nil {
			auth.log.LogError("Error while VerfiyOTP", err)
			return false, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		if !sts {
			auth.log.LogError("Invalid OTP")
			return false, status.Errorf(codes.PermissionDenied, "Invalid OTP")
		}
	} else {
		if payload.Value != phone {
			auth.log.LogError("Forbidden")
			return false, status.Errorf(codes.PermissionDenied, "Forbidden")
		}
		phoneNumber := fmt.Sprintf("+91%s", phone)
		sts, err := auth.twilioService.VerfiyOTP(phoneNumber, code)
		if err != nil {
			auth.log.LogError("Error while VerfiyOTP", err)
			return false, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		if !sts {
			auth.log.LogError("Invalid OTP")
			return false, status.Errorf(codes.PermissionDenied, "Invalid OTP")
		}
	}
	return true, nil
}

func (auth *authenticationService) SignUpUser(ctx context.Context, req *proto.SignupRequest) (*proto.Response, error) {

	err := helpers.ValidateEmailOrPhone(req.GetEmail(), req.GetPhone())
	if err != nil {
		auth.log.LogError("Error while ValidateEmailOrPhone", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Email or Phone")
	}

	payload := ctx.Value(utils.AuthorizationPayloadKey).(*token.SessionPayload)
	verify, err := auth.verifyOTP(payload, req.Email, req.Phone, req.Code, utils.SIGN_UP)

	if !verify {
		return nil, err
	}
	id, err := auth.storage.InsertUser(req)
	if err != nil {
		auth.log.LogError("Error while InsertUser", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	token, err := auth.token.CreateAccessToken(id, req.UserName, req.UserType)
	if err != nil {
		auth.log.LogError("Error while CreateAccessToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "User Signup Successfully",
	}

	mdOut := metadata.New(map[string]string{
		"access-token": token,
	})
	return response, grpc.SetHeader(ctx, mdOut)
}

func (auth *authenticationService) LoginUser(ctx context.Context, req *proto.LoginRequest) (*proto.Response, error) {
	res, err := auth.storage.CheckDataExist("user_name", req.GetUserName())
	if err != nil {
		auth.log.LogError("Error while CheckDataExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		auth.log.LogError("User doesn't exist", req.GetUserName())
		return nil, status.Errorf(codes.NotFound, "User doesn't exist")
	}

	res, err = auth.storage.CheckPassword(req.UserName, req.Password)
	if err != nil {
		auth.log.LogError("Error while CheckPassword", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		auth.log.LogError("Wrong password")
		return nil, status.Errorf(codes.PermissionDenied, "User crediantials doesn't match")
	}

	payload, err := auth.storage.GetTokenPayload("user_name", req.UserName)
	if err != nil {
		auth.log.LogError("Error while GetTokenPayload", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	token, err := auth.token.CreateAccessToken(payload.UserId, req.UserName, payload.UserType)
	if err != nil {
		auth.log.LogError("Error while CreateAccessToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "User loggedin Successfully",
	}

	mdOut := metadata.New(map[string]string{
		"access-token": token,
	})
	return response, grpc.SetHeader(ctx, mdOut)
}

func (auth *authenticationService) ForgotPassword(ctx context.Context, req *proto.ForgotPasswordRequest) (*proto.Response, error) {

	err := validateForgotPasswordRequest(req)
	if err != nil {
		auth.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	payload := ctx.Value(utils.AuthorizationPayloadKey).(*token.SessionPayload)
	verify, err := auth.verifyOTP(payload, req.Email, req.Phone, req.Code, utils.FORGOT_PASSWORD)
	if !verify {
		return nil, err
	}

	err = auth.storage.ChangePassword(req)
	if err != nil {
		auth.log.LogError("Error while ChangePassword", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "Password changed successfully",
	}
	return response, nil
}

func (auth *authenticationService) SSOAuth(ctx context.Context, req *proto.SsoRequest) (*proto.Response, error) {
	err := validateSsoRequest(req)
	if err != nil {
		auth.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	token := ctx.Value(utils.AuthorizationPayloadKey).(string)

	var email, fullname string

	switch req.Action {
	case utils.SSO_ANDROID:
		email, fullname, err = sso.GoogleOauth(token, auth.config.OAuth.AndroidClientId, auth.log)
		if err != nil {
			if err.Error() == "missing claims" {
				auth.log.LogError("conflict occurs, missing claims :", err)
				return nil, status.Errorf(codes.NotFound, "conflict occurs, missing claims")
			}
			auth.log.LogError("Error while GoogleOauth in sso-android", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	case utils.SSO_IOS:
		email, fullname, err = sso.GoogleOauth(token, auth.config.OAuth.IOSClientId, auth.log)
		if err != nil {
			if err.Error() == "missing claims" {
				auth.log.LogError("conflict occurs, missing claims :", err)
				return nil, status.Errorf(codes.NotFound, "conflict occurs, missing claims")
			}
			auth.log.LogError("Error while GoogleOauth in sso-ios", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	default:
		auth.log.LogError("Bad Requst", req)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	//checks if the user exist or not
	exist, err := auth.storage.CheckDataExist("email", email)
	if err != nil {
		auth.log.LogError("Error while CheckDataExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	//already exist so login
	if exist {
		payload, err := auth.storage.GetTokenPayload("email", email)
		if err != nil {
			auth.log.LogError("Error while GetTokenPayload", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}

		token, err := auth.token.CreateAccessToken(payload.UserId, payload.UserName, payload.UserType)
		if err != nil {
			auth.log.LogError("Error while CreateAccessToken", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}

		response := &proto.Response{
			Code:    200,
			Status:  true,
			Message: "User loggedin Successfully",
		}

		mdOut := metadata.New(map[string]string{
			"access-token": token,
		})
		return response, grpc.SetHeader(ctx, mdOut)
	} else {
		signupReq := &proto.SignupRequest{
			Email:    email,
			UserName: fullname,
			UserType: req.GetUserType(),
			Password: email + fullname + req.UserType,
		}
		id, err := auth.storage.InsertUser(signupReq)
		if err != nil {
			auth.log.LogError("Error while InsertUser", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}

		token, err := auth.token.CreateAccessToken(id, fullname, req.UserType)
		if err != nil {
			auth.log.LogError("Error while CreateAccessToken", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}

		response := &proto.Response{
			Code:    200,
			Status:  true,
			Message: "User Signup Successfully",
		}

		mdOut := metadata.New(map[string]string{
			"access-token": token,
		})
		return response, grpc.SetHeader(ctx, mdOut)
	}
}
