package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *productService) GetProducts(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductsResponse, error) {
	payload, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	var productRes []*entity.GetProductRes
	if payload.UserType == utils.USER {
		productRes, err = handler.storage.GetProductsForUser(req.MerchantId, payload.UserID)
		if err != nil {
			handler.log.LogError("Error while GetProducts", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	} else {
		productRes, err = handler.storage.GetProductsForMerchants(req.MerchantId)
		if err != nil {
			handler.log.LogError("Error while GetProducts", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}

	for _, product := range productRes {
		if product.ProductImages != nil {
			for i, image := range product.ProductImages {
				url, err := handler.s3.GetPreSignedURL(image)
				if err != nil {
					handler.log.LogError("Error while GetPreSignedURL", err)
					return nil, status.Errorf(codes.Internal, utils.InternalServerError)
				}
				product.ProductImages[i] = url
			}
		}
	}
	var products []*proto.ProductResponse
	for _, product := range productRes {
		newProduct := &proto.ProductResponse{
			Id:            product.ID,
			MerchantId:    product.MerchantID,
			ProductImages: product.ProductImages,
			ProductName:   product.ProductName,
			Description:   product.Description,
			CategoryName:  product.CategoryName,
			// Size:          *product.Size,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			ReviewStar:    product.ReviewStar,
			WishlistId:    product.WishlistID,
		}
		products = append(products, newProduct)
	}

	response := &proto.GetProductsResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Products fetched successfully",
		Data:    products,
	}
	return response, nil
}

func (handler *productService) GetProductById(ctx context.Context, req *proto.ProductIdRequest) (*proto.GetProductByIdResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var product *entity.GetProductRes
	var err error
	if payload.UserType == utils.USER {
		product, err = handler.storage.GetProductByIdForUser(req.GetProductId(), payload.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				handler.log.LogError("Error while GetProductById Not found", err)
				return nil, status.Errorf(codes.NotFound, "No found")
			}
			handler.log.LogError("Error while GetProductById", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	} else {
		product, err = handler.storage.GetProductByIdForMerchant(req.GetProductId())
		if err != nil {
			if err == sql.ErrNoRows {
				handler.log.LogError("Error while GetProductById Not found", err)
				return nil, status.Errorf(codes.NotFound, "No found")
			}
			handler.log.LogError("Error while GetProductById", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}

	if product.ProductImages != nil {
		for i, image := range product.ProductImages {
			url, err := handler.s3.GetPreSignedURL(image)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			product.ProductImages[i] = url
		}
	}

	productRes := &proto.ProductResponse{
		Id:            product.ID,
		MerchantId:    product.MerchantID,
		ProductImages: product.ProductImages,
		ProductName:   product.ProductName,
		Description:   product.Description,
		CategoryName:  product.CategoryName,
		Size:          *product.Size,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		ReviewStar:    product.ReviewStar,
		WishlistId:    product.WishlistID,
	}

	response := &proto.GetProductByIdResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product fetched successfully",
		Data:    productRes,
	}

	return response, nil
}
