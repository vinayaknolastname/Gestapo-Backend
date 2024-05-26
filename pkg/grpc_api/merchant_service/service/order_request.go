package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/db/entity"
	user_entity "github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (handler *merchantService) GetMerchantOrders(ctx context.Context, in *proto.GetOrdersRequest) (*proto.GetOrderResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve merchant payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &user_entity.GetOrdersReq{
		Type: in.GetType(),
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	serviceToken, err := handler.token.CreateServiceToken(payload.UserID, payload.UserType, "order")
	if err != nil {
		handler.log.LogError("error while generating service token in CreateOrder", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	conn, err := service_helper.ConnectEndpoints(handler.config.ServerAddress.Order.Address, "order", handler.log)
	if err != nil {
		handler.log.LogError("error while connecting order service :", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	defer conn.Close()

	orderClient := proto.NewOrderServiceClient(conn)
	serviceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	serviceCtx = metadata.NewOutgoingContext(serviceCtx, metadata.New(map[string]string{
		token.ServiceToken: fmt.Sprint(utils.AuthorizationTypeBearer, " ", serviceToken),
	}))
	defer cancel()

	response, err := orderClient.GetMerchantOrders(serviceCtx, in)
	if err != nil {
		handler.log.LogError("error parsing order service context :", err)
		return nil, err
	}
	return response, nil
}

func (handler *merchantService) UpdateOrderStatus(ctx context.Context, in *proto.UpdateOrderRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve merchant payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &entity.UpdateOrderReq{
		OrderItemID: in.GetOrderItemId(),
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	serviceToken, err := handler.token.CreateServiceToken(payload.UserID, payload.UserType, "order")
	if err != nil {
		handler.log.LogError("error while generating service token in CreateOrder", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	conn, err := service_helper.ConnectEndpoints(handler.config.ServerAddress.Order.Address, "order", handler.log)
	if err != nil {
		handler.log.LogError("error while connecting order service :", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	defer conn.Close()

	orderClient := proto.NewOrderServiceClient(conn)
	serviceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	serviceCtx = metadata.NewOutgoingContext(serviceCtx, metadata.New(map[string]string{
		token.ServiceToken: fmt.Sprint(utils.AuthorizationTypeBearer, " ", serviceToken),
	}))
	defer cancel()

	response, err := orderClient.UpdateOrderStatus(serviceCtx, in)
	if err != nil {
		handler.log.LogError("error parsing order service context :", err)
		return nil, err
	}
	return response, nil
}
