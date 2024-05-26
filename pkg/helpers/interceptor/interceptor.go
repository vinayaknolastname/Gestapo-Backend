package interceptor

import (
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
)

type Interceptor struct {
	token token.Maker
	log   logger.Logger
}

func NewInterceptor(token token.Maker, log logger.Logger) *Interceptor {
	return &Interceptor{
		token: token,
		log:   log,
	}
}
