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

func (handler *userService) AddAddress(ctx context.Context, in *proto.AddAddressRequest) (*proto.Response, error) {
	req := &entity.AddAddressReq{
		Title:       in.GetTitle(),
		AddressLine: in.GetAddressLine(),
		Country:     in.GetCountry(),
		City:        in.GetCity(),
		PostalCode:  in.PostalCode,
		Landmark:    in.Landmark,
		IsDefault:   in.IsDefault,
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	req.UserID = payload.UserID

	err = handler.storage.AddAddress(req)
	if err != nil {
		handler.log.LogError("Error while AddAddress", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Address added successfully",
	}
	return response, nil
}

func (handler *userService) GetAddresses(ctx context.Context, in *proto.Request) (*proto.GetAddressesResponse, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	addressEntities, err := handler.storage.GetAddresses(payload.UserID)
	if err != nil {
		handler.log.LogError("Error while GetAddresses", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	var addresses []*proto.AddressesResponse
	for _, address := range addressEntities {
		value := &proto.AddressesResponse{
			AddressId:   address.AddressID,
			Title:       address.Title,
			AddressLine: address.AddressLine,
		}
		addresses = append(addresses, value)
	}

	response := &proto.GetAddressesResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Address fetched successfully",
		Data:    addresses,
	}
	return response, nil
}
func (handler *userService) GetAddressByID(ctx context.Context, in *proto.AddressIdRequest) (*proto.GetAddressByIdResponse, error) {
	// ------------ Removing this condition because other service may use this function to get the
	// payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	// if !ok {
	// 	err := errors.New("unable to retrieve user payload from context")
	// 	handler.log.LogError("Error", err)
	// 	return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	// }

	address, err := handler.storage.GetAddressById(in.GetAddressId())
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetAddressById")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetAddressById", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}
	// ------------ Removing this condition because other service may use this function to get the
	// if payload.UserID != addressEntity.UserID {
	// 	handler.log.LogError("Address doesn't belong the user")
	// 	return nil, status.Errorf(codes.PermissionDenied, utils.PermissionDenied)
	// }

	addressRes := &proto.AddressesResponse{
		AddressId:   address.AddressID,
		Title:       address.Title,
		AddressLine: address.AddressLine,
		Country:     address.Country,
		City:        address.City,
		PostalCode:  address.PostalCode,
		Landmark:    address.Landmark,
		IsDefault:   address.IsDefault,
	}

	response := &proto.GetAddressByIdResponse{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Address fetched successfully",
		Data:    addressRes,
	}
	return response, nil
}

func (handler *userService) EditAddress(ctx context.Context, in *proto.EditAddressRequest) (*proto.Response, error) {
	req := &entity.EditAddressReq{
		Title:       in.Title,
		AddressLine: in.AddressLine,
		Country:     in.Country,
		City:        in.City,
		PostalCode:  in.PostalCode,
		Landmark:    in.Landmark,
		IsDefault:   in.IsDefault,
	}
	err := helpers.ValidateBody(nil, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		return nil, status.Errorf(codes.InvalidArgument, utils.InvalidRequest)
	}

	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	addressEntity, err := handler.storage.GetAddressById(in.GetAddressId())
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetAddressById")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetAddressById", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	if payload.UserID != addressEntity.UserID {
		handler.log.LogError("Address doesn't belong the user")
		return nil, status.Errorf(codes.PermissionDenied, utils.PermissionDenied)
	}

	err = handler.storage.EditAddress(in.GetAddressId(), req)
	if err != nil {
		handler.log.LogError("Error while EditAddress", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Address updated successfully",
	}

	return response, nil
}

func (handler *userService) DeleteAddress(ctx context.Context, in *proto.AddressIdRequest) (*proto.Response, error) {
	payload, ok := ctx.Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	addressEntity, err := handler.storage.GetAddressById(in.GetAddressId())
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetAddressById")
			return nil, status.Errorf(codes.NotFound, utils.NotFound)
		}
		handler.log.LogError("Error while GetAddressById", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	if payload.UserID != addressEntity.UserID {
		handler.log.LogError("Address doesn't belong the user")
		return nil, status.Errorf(codes.PermissionDenied, utils.PermissionDenied)
	}

	err = handler.storage.DeleteAddress(in.GetAddressId())
	if err != nil {
		handler.log.LogError("Error while DeleteAddress", err)
		return nil, status.Errorf(codes.Internal, utils.InternalServerError)
	}

	response := &proto.Response{
		Code:    http.StatusOK,
		Status:  true,
		Message: "Address deleted successfully",
	}

	return response, nil
}
