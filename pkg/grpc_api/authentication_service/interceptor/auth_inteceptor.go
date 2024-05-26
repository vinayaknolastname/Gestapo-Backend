package interceptor

import (
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
)

type AuthInterceptor struct {
	token token.Maker
	log   logger.Logger
}

func NewAuthInterceptor(token token.Maker, log logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		token: token,
		log:   log,
	}
}
