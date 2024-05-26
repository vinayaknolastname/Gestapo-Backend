package service

import (
	"context"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (admin *adminService) GetUsers(ctx context.Context, in *proto.Request) (*proto.GetUsersResponse, error) {
	userEntites, err := admin.storage.GetUsers()
	if err != nil {
		admin.log.LogError("Error while GetUsers", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var users []*proto.UserResponse
	for _, user := range userEntites {
		var dob *timestamppb.Timestamp
		if user.DOB != nil {
			dob = timestamppb.New(*user.DOB)
		}
		userRes := &proto.UserResponse{
			Id:           user.ID,
			ProfileImage: user.ProfileImage,
			FullName:     user.FullName,
			UserName:     user.UserName,
			Phone:        user.Phone,
			Email:        user.Email,
			Dob:          dob,
			Gender:       user.Gender,
			UserType:     user.UserType,
		}
		users = append(users, userRes)
	}

	response := &proto.GetUsersResponse{
		Code:    200,
		Status:  true,
		Message: "Users fetched successfully",
		Data:    users,
	}
	return response, nil
}
