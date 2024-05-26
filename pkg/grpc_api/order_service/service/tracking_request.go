package service

import (
	"context"
	"net/http"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/order_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *orderService) GetOrderTrackingDetails(ctx context.Context, in *proto.GetTrackingDetailsRequest) (*proto.GetTrackingDetailsResponse, error) {
	req := &entity.GetTrackingDetailsReq{
		OrderItemID: in.GetOrderItemId(),
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	trackingEntities, err := handler.storage.GetOrderTrackingDetails(in.GetOrderItemId())
	if err != nil {
		handler.log.LogError("Error while GetOrderTrackingDetails", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var trackingItems []*proto.TrackingItemsResponse
	var status int
	for _, item := range trackingEntities {
		status = int(item.Status)
		val := &proto.TrackingItemsResponse{
			Title:   item.Title,
			Summary: item.Summary,
			Time:    timestamppb.New(item.Time),
		}
		trackingItems = append(trackingItems, val)
	}

	response := &proto.GetTrackingDetailsResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Tracking details fetched successfully",
		Data: &proto.TrackingDetailsResponse{
			Status:  int32(status),
			Details: trackingItems,
		},
	}

	return response, nil
}
