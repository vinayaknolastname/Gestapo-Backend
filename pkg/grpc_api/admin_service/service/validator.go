package service

import (
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateAddCategoryRequest(req *proto.AddCategoryRequest) error {
	// Check if either email or phone is provided
	if helpers.IsEmpty(req.GetCategoryName()) {
		return status.Errorf(codes.InvalidArgument, "Category name must not be empty")

	}
	return nil

}
