package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	signUp         string = "/pb.AuthenticationService/SignUpUser"
	forgotPassword string = "/pb.AuthenticationService/ForgotPassword"
	sso            string = "/pb.AuthenticationService/SSOAuth"
)

// AuthMiddleware is a gRPC unary server interceptor for authentication.
func (interceptor *AuthInterceptor) AuthMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		interceptor.log.LogInfo("Calling gRPC meathod :", info.FullMethod)
		if ok := isAuthenticationNeeded(info.FullMethod); !ok {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err := errors.New("metadata is not provided")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		authorizationHeaders := md.Get(utils.AuthorizationKey)
		if len(authorizationHeaders) == 0 {
			err := errors.New("authorization header is not provided")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		authorizationHeader := authorizationHeaders[0]
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != utils.AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		token := fields[1]
		// Verify and parse the token
		payload, err := interceptor.token.VerifySessionToken(token)
		if err != nil {
			err := fmt.Errorf("error while VerifySessionToken: %s", err.Error())
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		// Add the payload to the context
		ctx = context.WithValue(ctx, utils.AuthorizationPayloadKey, payload)

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) AuthSsoMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ok := isSSONeeded(info.FullMethod); !ok {
			interceptor.log.LogInfo("Calling gRPC meathod :", info.FullMethod)
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err := errors.New("metadata is not provided")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
		authorizationHeaders := md.Get(utils.AuthorizationKey)
		if len(authorizationHeaders) == 0 {
			err := errors.New("authorization header is not provided")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())

		}
		authorizationHeader := authorizationHeaders[0]
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())

		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != utils.AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			interceptor.log.LogError("Error : ", err)
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		token := fields[1]

		ctx = context.WithValue(ctx, utils.AuthorizationPayloadKey, token)
		return handler(ctx, req)
	}
}

// isAuthenticationNeeded returns true if the route needed authentication middleware
func isAuthenticationNeeded(route string) bool {
	switch route {
	case signUp, forgotPassword:
		return true
	}
	return false
}

// isAuthenticationNeeded returns true if the route needed authentication middleware
func isSSONeeded(route string) bool {
	return route == sso
}
