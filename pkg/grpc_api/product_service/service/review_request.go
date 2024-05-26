package service

import (
	"context"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *productService) AddProductReview(ctx context.Context, in *proto.AddReviewRequest) (*proto.Response, error) {
	paylaod, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &entity.AddReviewReq{
		ProductID:   in.GetProductId(),
		OrderItemID: in.GetOrderItemId(),
		UserID:      paylaod.UserID,
		Star:        in.GetStart(),
		Review:      in.GetReview(),
	}
	err = helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	res, err := handler.storage.IsUserCanAddReview(req.OrderItemID, req.UserID)
	if err != nil {
		handler.log.LogError("Error while IsUserCanAddReview", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		handler.log.LogError("User Cannot add review")
		return nil, status.Errorf(codes.PermissionDenied, "User cannot add review")
	}

	res, err = handler.storage.IsUserAlreadyAddedReview(req.ProductID, req.UserID)
	if err != nil {
		handler.log.LogError("Error while IsUserAlreadyAddedReview", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if res {
		handler.log.LogError("User already added review")
		return nil, status.Errorf(codes.PermissionDenied, "User already added review")
	}

	err = handler.storage.AddProductReview(req)
	if err != nil {
		handler.log.LogError("Error while AddProductReview", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product Review added successfully",
	}
	return response, nil
}
