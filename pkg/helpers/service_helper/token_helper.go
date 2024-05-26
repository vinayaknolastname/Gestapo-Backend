package service_helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ValidateServiceToken(ctx context.Context, log logger.Logger, tokenMaker token.Maker) (*token.ServicePayload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := errors.New("metadata is not provided")
		log.LogError(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	authorizationHeaders := md.Get(token.ServiceToken)
	if len(authorizationHeaders) == 0 {
		err := errors.New("authorization header is not provided")
		log.LogError(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	authorizationHeader := authorizationHeaders[0]
	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("invalid authorization header format")
		log.LogError(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != utils.AuthorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
		log.LogError(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	token := fields[1]
	// Verify and parse the token
	payload, err := tokenMaker.VerifyServiceToken(token)
	if err != nil {
		err := fmt.Errorf("error while VerifySessionToken: %s", err.Error())
		log.LogError(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return payload, nil

}
