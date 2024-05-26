package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (handler *userService) AddProductReview(ctx context.Context, in *proto.AddReviewRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	res, err := handler.storage.CheckDataExist("products", "id", in.GetProductId())
	if err != nil {
		handler.log.LogError("Error while CheckDataExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		handler.log.LogError("product not found")
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	res, err = handler.storage.CheckDataExist("order_items", "id", in.GetOrderItemId())
	if err != nil {
		handler.log.LogError("Error while CheckDataExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		handler.log.LogError("order item not found")
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	serviceToken, err := handler.token.CreateServiceToken(payload.UserID, payload.UserType, "product")
	if err != nil {
		handler.log.LogError("error while generating service token in AddProductReview", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	conn, err := service_helper.ConnectEndpoints(handler.config.ServerAddress.Product.Address, "product", handler.log)
	if err != nil {
		handler.log.LogError("error while connecting order service :", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	defer conn.Close()

	productClient := proto.NewProductServiceClient(conn)
	serviceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	serviceCtx = metadata.NewOutgoingContext(serviceCtx, metadata.New(map[string]string{
		token.ServiceToken: fmt.Sprint(utils.AuthorizationTypeBearer, " ", serviceToken),
	}))
	defer cancel()

	response, err := productClient.AddProductReview(serviceCtx, in)
	if err != nil {
		handler.log.LogError("error parsing product service context :", err)
		return nil, err
	}
	return response, nil
}
