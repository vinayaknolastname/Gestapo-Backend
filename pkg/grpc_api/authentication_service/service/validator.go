package service

import (
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateSendOTPRequest(req *proto.SendOTPRequest) error {
	// Check if either email or phone is provided
	err := helpers.ValidateEmailOrPhone(req.GetEmail(), req.GetPhone())
	if err != nil {
		return err
	}
	// Validate action
	if !utils.IsSupportedSignupAction(req.GetAction()) {
		return status.Errorf(codes.InvalidArgument, "action must be either 'signup' or 'forget'")
	}
	return nil
}

func validateForgotPasswordRequest(req *proto.ForgotPasswordRequest) error {
	// Check if either email or phone is provided
	err := helpers.ValidateEmailOrPhone(req.GetEmail(), req.GetPhone())
	if err != nil {
		return err
	}

	if !utils.IsValidPassword(req.GetPassword()) {
		return status.Errorf(codes.InvalidArgument, "password must greater than 6 and less than 100")
	}

	if !utils.IsValidPassword(req.GetPassword()) {
		return status.Errorf(codes.InvalidArgument, "code should be 6 digit")
	}
	return nil
}

func validateSsoRequest(req *proto.SsoRequest) error {

	if !utils.IsSupportedSSOAction(req.GetAction()) {
		return status.Errorf(codes.InvalidArgument, "invalid action")
	}

	if !utils.IsSupportedUsers(req.GetUserType()) {
		return status.Errorf(codes.InvalidArgument, "invalid user type")
	}
	return nil
}
