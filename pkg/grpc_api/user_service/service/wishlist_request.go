package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *userService) AddRemoveWishlist(ctx context.Context, in *proto.AddRemoveWishlistRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &entity.AddRemoveWishlistReq{
		Action:    in.GetAction(),
		ProductID: in.GetProductId(),
		UserID:    payload.UserID,
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	res, err := handler.storage.CheckDataExist("products", "id", req.ProductID)
	if err != nil {
		handler.log.LogError("Error while CheckUserExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		err = fmt.Errorf("product does'nt exist using %s", req.ProductID)
		handler.log.LogError(err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	switch req.Action {
	case utils.ADD_WISHLIST:
		return handler.addToWishlist(req)
	case utils.REMOVE_WISHLIST:
		return handler.removeFromWishlist(req)
	default:
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}
}

func (handler *userService) addToWishlist(req *entity.AddRemoveWishlistReq) (*proto.Response, error) {
	res, err := handler.storage.AlreadyInWishlist(req)
	if err != nil {
		handler.log.LogError("Error while AlreadyInWishlist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if res {
		handler.log.LogError("Error while AlreadyInWishlist", err)
		return nil, status.Errorf(codes.AlreadyExists, utils.AlreadyExists)
	}

	err = handler.storage.AddToWishlist(req)
	if err != nil {
		handler.log.LogError("Error while AddToWishlist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product added to wishlist",
	}

	return response, nil
}
func (handler *userService) removeFromWishlist(req *entity.AddRemoveWishlistReq) (*proto.Response, error) {
	res, err := handler.storage.AlreadyInWishlist(req)
	if err != nil {
		handler.log.LogError("Error while AlreadyInWishlist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		handler.log.LogError("Error while AlreadyInWishlist", err)
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	err = handler.storage.RemoveFromWishlist(req)
	if err != nil {
		handler.log.LogError("Error while RemoveFromWishlist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product removed from wishlist",
	}

	return response, nil
}

func (handler *userService) GetWishlist(ctx context.Context, in *proto.Request) (*proto.GetWishlistResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	productEntities, err := handler.storage.GetWishlistProducts(payload.UserID)
	if err != nil {
		handler.log.LogError("Error while GetWishlistProducts", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, product := range productEntities {
		if product.ProductImages[0] != "" {
			url, err := handler.s3.GetPreSignedURL(product.ProductImages[0])
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL for product.ProductImages[0]", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			product.ProductImages[0] = url
		}
	}

	var products []*proto.ProductResponse
	for _, product := range productEntities {
		newProduct := &proto.ProductResponse{
			Id:            product.ID,
			ProductImages: product.ProductImages,
			ProductName:   product.ProductName,
			Price:         product.Price,
		}
		products = append(products, newProduct)
	}

	response := &proto.GetWishlistResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Wishlist fetched successfully",
		Data:    products,
	}

	return response, nil
}
