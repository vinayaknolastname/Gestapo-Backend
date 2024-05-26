package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (handler *merchantService) GetProducts(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductsResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve merchant payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	serviceToken, err := handler.token.CreateServiceToken(payload.UserID, payload.UserType, "product")
	if err != nil {
		handler.log.LogError("error while generating service token in GetProducts", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	conn, err := service_helper.ConnectEndpoints(handler.config.ServerAddress.Product.Address, "product", handler.log)
	if err != nil {
		handler.log.LogError("error while connecting product service :", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	defer conn.Close()

	productClient := proto.NewProductServiceClient(conn)
	serviceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	serviceCtx = metadata.NewOutgoingContext(serviceCtx, metadata.New(map[string]string{
		token.ServiceToken: fmt.Sprint(utils.AuthorizationTypeBearer, " ", serviceToken),
	}))
	defer cancel()

	response, err := productClient.GetProducts(serviceCtx, &proto.GetProductRequest{MerchantId: req.MerchantId})
	if err != nil {
		handler.log.LogError("error parsing product service context :", err)
		return nil, err
	}
	return response, nil
}

func (handler *merchantService) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve merchant payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	serviceToken, err := handler.token.CreateServiceToken(payload.UserID, payload.UserType, "product")
	if err != nil {
		handler.log.LogError("error while generating service token in DeleteProduct", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	conn, err := service_helper.ConnectEndpoints(handler.config.ServerAddress.Product.Address, "product", handler.log)
	if err != nil {
		handler.log.LogError("error while connecting product service :", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	defer conn.Close()

	productClient := proto.NewProductServiceClient(conn)
	serviceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	serviceCtx = metadata.NewOutgoingContext(serviceCtx, metadata.New(map[string]string{
		token.ServiceToken: fmt.Sprint(utils.AuthorizationTypeBearer, " ", serviceToken),
	}))
	defer cancel()

	productRes, err := productClient.GetProductById(serviceCtx, &proto.ProductIdRequest{
		ProductId: req.GetProductId(),
	})
	if err != nil {
		handler.log.LogError("Error while retrieving product", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !productRes.Status {
		if productRes.Code == int32(codes.NotFound) {
			handler.log.LogError("Error while GetProductById product Not found")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetProductById")
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	if productRes.Data.MerchantId != nil && *productRes.Data.MerchantId != payload.UserID {
		handler.log.LogError("unauthorized: product does not belong to the authenticated merchant")
		return nil, status.Errorf(codes.PermissionDenied, "product does not belong to the authenticated merchant")
	}

	err = handler.storage.DeleteProduct(req.GetProductId())
	if err != nil {
		handler.log.LogError("Error while DeleteProduct", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	for _, key := range productRes.Data.ProductImages {
		err := handler.s3.DeleteKey(key)
		if err != nil {
			handler.log.LogError("Error deleting file from S3", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product deleted succesfully",
	}
	return response, nil
}
