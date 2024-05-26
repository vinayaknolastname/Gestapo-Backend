package interceptor

import (
	"context"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator interface {
	Validate() error
}

func (interceptor *AuthInterceptor) AuthValidator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		t := reflect.TypeOf(req)
		interceptor.log.LogInfo("type ", t)
		if r, ok := req.(Validator); ok {
			if err := r.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}
		return handler(ctx, req)
	}
}

// func (req *proto.SendOTPRequest) Validate() error {
// 	// Check if either email or phone is provided
// 	err := helpers.ValidateEmailOrPhone(req.GetEmail(), req.GetPhone())
// 	if err != nil {
// 		return err
// 	}
// 	// Validate action
// 	if utils.IsSupportedSignupAction(req.GetAction()) {
// 		return status.Errorf(codes.InvalidArgument, "action must be either 'signup' or 'forget'")
// 	}

// 	return nil
// }
