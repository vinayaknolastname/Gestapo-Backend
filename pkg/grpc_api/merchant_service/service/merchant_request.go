package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *merchantService) GetProfile(ctx context.Context, req *proto.GetMerchantProfileRequest) (*proto.GetMerchantProfileResponse, error) {
	if req.GetUserId() == "" {
		handler.log.LogError("Error while Getting user id")
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}
	res, err := handler.storage.CheckDataExist("user_data", "id", req.GetUserId())
	if err != nil {
		handler.log.LogError("Error while CheckUserExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		err = fmt.Errorf("account does'nt exist using %s", req.GetUserId())
		handler.log.LogError(err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	userData, err := handler.storage.GetProfile(req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetProfile", err)
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetProfile", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	if userData.ProfileImage != nil && *userData.ProfileImage != "" {
		url, err := handler.s3.GetPreSignedURL(*userData.ProfileImage)
		if err != nil {
			handler.log.LogError("Error while GetPreSignedURL", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		userData.ProfileImage = &url
	}

	var dob *timestamppb.Timestamp
	if userData.DOB != nil {
		dob = timestamppb.New(*userData.DOB)
	}

	merchantData := &proto.MerchantResponse{
		Id:           userData.ID,
		ProfileImage: userData.ProfileImage,
		FullName:     userData.FullName,
		UserName:     userData.UserName,
		Phone:        userData.Phone,
		Email:        userData.Email,
		Dob:          dob,
		Gender:       userData.Gender,
		UserType:     userData.UserType,
	}
	respone := &proto.GetMerchantProfileResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Profile fetched sucessfull",
		Data:    merchantData,
	}
	return respone, nil
}
