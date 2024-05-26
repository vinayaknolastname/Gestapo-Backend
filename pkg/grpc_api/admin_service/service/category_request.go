package service

import (
	"context"
	"fmt"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *adminService) CreateCategory(ctx context.Context, req *proto.AddCategoryRequest) (*proto.Response, error) {
	err := validateAddCategoryRequest(req)
	if err != nil {
		handler.log.LogError("Error while validateAddCategoryRequest", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	res, err := handler.storage.CheckCategoryExist(req.GetCategoryName())
	if err != nil {
		handler.log.LogError("Error while CheckCategoryExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if res {
		err = fmt.Errorf("category already exist: %s", req.GetCategoryName())
		handler.log.LogError("Error ", err)
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	err = handler.storage.AddCategory(req)
	if err != nil {
		handler.log.LogError("Error while InsertCategory", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "Category insterted Successfully",
	}
	return response, err
}

func (handler *adminService) GetCategories(ctx context.Context, in *proto.Request) (*proto.GetCategoryResponse, error) {
	res, err := handler.storage.GetCategories()
	if err != nil {
		handler.log.LogError("Error while GetCategories", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.GetCategoryResponse{
		Code:    200,
		Status:  true,
		Message: "Categories fetched successfull",
		Data:    res,
	}
	return response, nil
}
