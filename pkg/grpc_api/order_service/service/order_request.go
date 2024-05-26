package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/order_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *orderService) CreateOrder(ctx context.Context, in *proto.CreateOrderRequest) (*proto.Response, error) {
	servicePayload, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	req := &entity.CreateOrderReq{
		AddressID:     in.GetAddressId(),
		CartID:        in.GetCartId(),
		PromoID:       in.PromoId,
		Amount:        float64(in.GetAmount()),
		PaymentMode:   in.GetPaymentMode(),
		TransactionID: in.TransactionId,
	}

	err = helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	if req.PaymentMode == utils.COD {
		res, err := handler.storage.CheckCODIsAvailable(servicePayload.UserID)
		if err != nil {
			handler.log.LogError("Error while CheckCODIsAvailable", err)
			return nil, status.Errorf(codes.Internal, utils.InternalServerError)
		}
		if !res {
			response := &proto.Response{
				Code:    http.StatusOK,
				Status:  false,
				Message: "User has to complete atleast 2 Order to avail COD",
			}
			return response, nil
		}
	}

	//assiging user id to create order request
	req.UserID = servicePayload.UserID
	err = handler.storage.CreateOrder(req)
	if err != nil {
		handler.log.LogError("Error while CreateOrder", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Orders created successfully",
	}
	return response, nil
}

func (handler *orderService) GetUserOrders(ctx context.Context, in *proto.GetOrdersRequest) (*proto.GetOrderResponse, error) {
	servicePayload, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	userOrdersEntities, err := handler.storage.GetUserOrders(servicePayload.UserID, in.Type)
	if err != nil {
		handler.log.LogError("Error while GetUserOrders", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, order := range userOrdersEntities {
		if order.ProductImage != "" {
			url, err := handler.s3.GetPreSignedURL(order.ProductImage)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			order.ProductImage = url
		}
	}

	var orders []*proto.OrderResponse
	for _, order := range userOrdersEntities {
		newProduct := &proto.OrderResponse{
			Id:           order.ID,
			ProductId:    order.ProductID,
			ProductImage: order.ProductImage,
			ProductName:  order.ProductName,
			Size:         float64(order.Size),
			Price:        order.Price,
			Status:       order.Status,
		}
		orders = append(orders, newProduct)
	}

	response := &proto.GetOrderResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Orders fetched successfully",
		Data:    orders,
	}
	return response, nil
}

func (handler *orderService) GetMerchantOrders(ctx context.Context, in *proto.GetOrdersRequest) (*proto.GetOrderResponse, error) {
	servicePayload, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	userOrdersEntities, err := handler.storage.GetMerchantOrders(servicePayload.UserID, in.Type)
	if err != nil {
		handler.log.LogError("Error while GetUserOrders", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	for _, order := range userOrdersEntities {
		if order.ProductImage != "" {
			url, err := handler.s3.GetPreSignedURL(order.ProductImage)
			if err != nil {
				handler.log.LogError("Error while GetPreSignedURL", err)
				return nil, status.Errorf(codes.Internal, utils.InternalServerError)
			}
			order.ProductImage = url
		}
	}

	var orders []*proto.OrderResponse
	for _, order := range userOrdersEntities {
		newProduct := &proto.OrderResponse{
			Id:           order.ID,
			ProductImage: order.ProductImage,
			ProductName:  order.ProductName,
			Size:         float64(order.Size),
			Price:        order.Price,
			Status:       order.Status,
		}
		orders = append(orders, newProduct)
	}

	response := &proto.GetOrderResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Orders fetched successfully",
		Data:    orders,
	}
	return response, nil
}

func (handler *orderService) UpdateOrderStatus(ctx context.Context, in *proto.UpdateOrderRequest) (*proto.Response, error) {
	payload, err := service_helper.ValidateServiceToken(ctx, handler.log, handler.token)
	if err != nil {
		handler.log.LogError("Error while ValidateServiceToken", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	res, err := handler.storage.IsMerchantCanUpdate(in.GetOrderItemId(), payload.UserID)
	if err != nil {
		handler.log.LogError("Error while CanEditDeleteCartItem", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if !res {
		err := errors.New("error while IsMerchantCanUpdate: Not found")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.NotFound, utils.NotFound)
	}

	statusCount, err := handler.storage.GetMerchantTrackingStatus(in.GetOrderItemId())
	if err != nil {
		handler.log.LogError("Error while GetMerchantTrackingStatus", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if statusCount >= 3 {
		err := errors.New("error while GetMerchantTrackingStatus: Oder is already completed")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.PermissionDenied, "Order is already completed")
	}

	err = handler.storage.UpdateOrderStatus(in.GetOrderItemId())
	if err != nil {
		handler.log.LogError("Error while UpdateOrderStatus", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Orders Status updated successfully",
	}
	return response, nil
}
