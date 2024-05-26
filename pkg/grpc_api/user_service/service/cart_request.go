package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *userService) AddProductToCart(ctx context.Context, in *proto.AddToCartRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &entity.AddToCartReq{
		ProductID: in.GetProductId(),
		Size:      float64(in.GetSize()),
		Quantity:  in.GetQuantity(),
		Price:     float64(in.GetPrice()),
		UserID:    payload.UserID,
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	res, err := handler.storage.CheckDataExist("carts", "user_id", req.UserID)
	if err != nil {
		handler.log.LogError("Error while CheckDataExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		err := handler.storage.CreateUserCart(req)
		if err != nil {
			handler.log.LogError("Error while CreateUserCart", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
	}

	err = handler.storage.AddToCard(req)
	if err != nil {
		handler.log.LogError("Error while AddToCard", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Product added to cart successfully",
	}
	return response, nil
}

func (handler *userService) GetCartItmes(ctx context.Context, in *proto.Request) (*proto.GetCartItemsResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	cartItemEntities, err := handler.storage.GetCartItems(payload.UserID)
	if err != nil {
		handler.log.LogError("Error while GetCartItems", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, product := range cartItemEntities {

		if product.ImageURL != "" {
			url, err := handler.s3.GetPreSignedURL(product.ImageURL)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			product.ImageURL = url
		}
	}

	var cartItems []*proto.CartItemResponse
	for _, cartItem := range cartItemEntities {
		item := &proto.CartItemResponse{
			CartId:            cartItem.CartID,
			CartItemId:        cartItem.CartItemID,
			ProductImage:      cartItem.ImageURL,
			Name:              cartItem.Name,
			Size:              float32(cartItem.Size),
			Price:             float32(cartItem.Price),
			Quantity:          cartItem.Quantity,
			AvailableQuantity: cartItem.AvailableQuantity,
		}
		cartItems = append(cartItems, item)
	}
	response := &proto.GetCartItemsResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Cart fetched successfully",
		Data:    cartItems,
	}
	return response, nil
}

func (handler *userService) CheckoutCartItems(ctx context.Context, in *proto.CheckoutCartItemsRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var cartItems []*entity.CheckoutReq
	for _, data := range in.Data {
		cartItem := entity.CheckoutReq{
			CartItemID: data.CartItemId,
			Quantity:   data.Quantity,
		}
		cartItems = append(cartItems, &cartItem)
	}

	req := &entity.CheckoutCartItemsReq{
		CartID: in.CartId,
		Data:   cartItems,
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	cartEntity, err := handler.storage.GetCartById(in.GetCartId())
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetCartById Not Found")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetCartById", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	if cartEntity.UserID != payload.UserID {
		handler.log.LogError("Cart doesn't belong the user")
		return nil, status.Errorf(codes.PermissionDenied, utils.PermissionDenied)
	}

	err = handler.storage.CheckoutCartItems(cartItems)
	if err != nil {
		handler.log.LogError("Error while CheckoutCartItems", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Cart Items updated successfully and proceed to checkout",
	}
	return response, nil
}

func (handler *userService) RemoveProductFromCart(ctx context.Context, in *proto.RemoveFromCartRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	//check cart_item is present or not
	res, err := handler.storage.CheckDataExist("cart_items", "id", in.GetCartItemId())
	if err != nil {
		handler.log.LogError("Error ", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		handler.log.LogError("Error CartItem not found")
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	res, err = handler.storage.CanDeleteCartItem(in.GetCartItemId(), payload.UserID)
	if err != nil {
		handler.log.LogError("Error while CanEditDeleteCartItem", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		err := errors.New("error while CanEditDeleteCartItem: Not found")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	err = handler.storage.RemoveFromCart(in.GetCartItemId(), payload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while RemoveFromCart")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		if err.Error() == "not deleted" {
			handler.log.LogError("Error while RemoveFromCart")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while RemoveFromCart", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Item deleted from cart successfully",
	}
	return response, nil
}
