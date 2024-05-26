package service

import (
	"context"
	"database/sql"
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

func (handler *userService) GetHome(ctx context.Context, req *proto.Request) (*proto.GetHomeResponse, error) {
	discount, err := handler.storage.GetDiscount()
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetDiscount No discouts", err)
			discount = nil
		} else {
			handler.log.LogError("Error while GetDiscount", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}
	//Converting the key to presigned url

	if discount != nil {
		url, err := handler.s3.GetPreSignedURL(discount.ProductImage)
		if err != nil {
			handler.log.LogError("Error while GetPreSignedURL product image", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		discount.ProductImage = url
	}

	merchantEntities, err := handler.storage.GetMerchants()
	if err != nil {
		handler.log.LogError("Error while GetMerchants", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, merchant := range merchantEntities {
		if merchant.ImageURL != nil {
			url, err := handler.s3.GetPreSignedURL(*merchant.ImageURL)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL for merchant.ImageURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			merchant.ImageURL = &url
		}
	}
	//For getting products from product service
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

	getProductsRes, err := productClient.GetProducts(serviceCtx, &proto.GetProductRequest{MerchantId: nil})
	if err != nil {
		handler.log.LogError("Error while GetProducts", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, product := range getProductsRes.Data {
		for i, image := range product.ProductImages {
			url, err := handler.s3.GetPreSignedURL(image)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			product.ProductImages[i] = url
		}

	}
	var discountRes *proto.DiscountResponse
	if discount != nil {
		discountRes = &proto.DiscountResponse{
			ProductId:    discount.ProductID,
			Name:         discount.Name,
			Description:  discount.Description,
			Percentage:   float32(discount.Percentage),
			ProductImage: discount.ProductImage,
			CardColor:    discount.CardColor,
		}
	}

	var merchants []*proto.UserResponse
	for _, merchantEntity := range merchantEntities {
		merchant := &proto.UserResponse{
			Id:           merchantEntity.MerchantID,
			ProfileImage: merchantEntity.ImageURL,
			FullName:     &merchantEntity.Name,
		}
		merchants = append(merchants, merchant)
	}

	home := &proto.HomeResponse{
		Discount:  discountRes,
		Merchants: merchants,
		Products:  getProductsRes.Data,
	}
	response := &proto.GetHomeResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Home data fetched successfully",
		Data:    home,
	}

	return response, nil
}
