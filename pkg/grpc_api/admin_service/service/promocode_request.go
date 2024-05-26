package service

import (
	"context"
	"fmt"

	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/admin_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *adminService) CreatePromocode(ctx context.Context, in *proto.CreatePromocodeRequest) (*proto.Response, error) {
	req := &entity.AddPromocodeReq{
		Code:        in.GetCode(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Percentage:  float64(in.GetPercentage()),
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	res, err := handler.storage.CheckPromocodeExist(req.Code)
	if err != nil {
		handler.log.LogError("Error while CheckPromocodeExist", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	if res {
		err = fmt.Errorf("category already exist: %s", req.Code)
		handler.log.LogError("Error ", err)
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	err = handler.storage.AddPromocode(req)
	if err != nil {
		handler.log.LogError("Error while InsertCategory", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	response := &proto.Response{
		Code:    200,
		Status:  true,
		Message: "Promocode insterted Successfully",
	}
	return response, err
}

func (handler *adminService) GetPromocodes(ctx context.Context, in *proto.Request) (*proto.GetPromocodeResponse, error) {
	promoEntities, err := handler.storage.GetPromocodes()
	if err != nil {
		handler.log.LogError("Error while GetCategories", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var promoCodes []*proto.PromocodeResponse
	for _, promo := range promoEntities {
		res := &proto.PromocodeResponse{
			PromoId:     promo.ID,
			Code:        promo.Code,
			Title:       promo.Title,
			Discription: promo.Description,
			Percentage:  float32(promo.Percentage),
		}

		promoCodes = append(promoCodes, res)

	}
	response := &proto.GetPromocodeResponse{
		Code:    200,
		Status:  true,
		Message: "Prmocodes fetched successfull",
		Data:    promoCodes,
	}
	return response, nil
}
