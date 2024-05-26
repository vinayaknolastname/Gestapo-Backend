package grpc_gateway

import (
	"context"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func newGateway(ctx context.Context, log logger.Logger, config config.Config, opts ...runtime.ServeMuxOption) (*runtime.ServeMux, error) {
	muxOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	gMux := runtime.NewServeMux(muxOption)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	//---------------Registering endpoints---------------------
	errAuthentication := registerAuthServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errAuthentication != nil {
		return nil, errAuthentication
	}

	errAdmin := registerAdminServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errAuthentication != nil {
		return nil, errAdmin
	}

	errUser := registerUserServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errUser != nil {
		return nil, errUser
	}

	errMerchant := registerMerchantServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errMerchant != nil {
		return nil, errMerchant
	}

	errProduct := registerProductServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errProduct != nil {
		return nil, errProduct
	}

	errOrder := registerOrderServiceEndPoints(ctx, log, config, gMux, dialOpts)
	if errOrder != nil {
		return nil, errOrder
	}

	return gMux, nil
}

func registerAuthServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.Authentication.Address
		err := proto.RegisterAuthenticationServiceHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering authentication endpoint.", err)
			return err
		}
	}
	return nil
}

func registerAdminServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.Admin.Address
		err := proto.RegisterAdminServiceHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering admin endpoint.", err)
			return err
		}
	}
	return nil
}

func registerUserServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.User.Address
		err := proto.RegisterUserServieHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering user endpoint.", err)
			return err
		}
	}
	return nil
}

func registerMerchantServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.Merchant.Address
		err := proto.RegisterMerchantServiceHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering merchant endpoint.", err)
			return err
		}
	}
	return nil
}

func registerProductServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.Product.Address
		err := proto.RegisterProductServiceHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering product endpoint.", err)
			return err
		}
	}
	return nil
}

func registerOrderServiceEndPoints(ctx context.Context, log logger.Logger, config config.Config, gMux *runtime.ServeMux, dialOpts []grpc.DialOption) error {
	var endpoint *string
	if config.ServerAddress != nil {
		endpoint = &config.ServerAddress.Order.Address
		err := proto.RegisterOrderServiceHandlerFromEndpoint(ctx, gMux, *endpoint, dialOpts)
		if err != nil {
			log.LogError("error in registering order endpoint.", err)
			return err
		}
	}
	return nil
}
