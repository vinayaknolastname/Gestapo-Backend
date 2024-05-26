package service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/authentication_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/service/cache"
	"github.com/akmal4410/gestapo/pkg/service/mail"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
	"github.com/akmal4410/gestapo/pkg/service/twilio"
)

// authenticationService serves gRPC requests for our e-commerce service.
type authenticationService struct {
	proto.UnimplementedAuthenticationServiceServer
	config        *config.Config
	log           logger.Logger
	s3            *s3.S3Service
	twilioService twilio.TwilioService
	emailService  mail.EmailService
	storage       *db.AuthStore
	token         token.Maker
	redis         cache.Cache
}

// NewAuthenticationService creates a new gRPC server.
func NewAuthenticationService(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *authenticationService {
	server := &authenticationService{
		config: config,
		log:    log,
		token:  tokenMaker,
	}
	s3 := s3.NewS3Service(
		config.AwsS3.BucketName,
		config.AwsS3.Region,
		config.AwsS3.AccessKey,
		config.AwsS3.SecretKey,
	)
	twilio := twilio.NewOTPService(server.config.Twilio)
	email := mail.NewGmailService(server.config.Email)
	authStore := db.NewAuthStore(storage)

	redis, err := cache.NewRedisCache(server.config.Redis)
	if err != nil {
		server.log.LogFatal("Error while Initializing NewRedisCache ", err)
	}

	server.twilioService = twilio
	server.s3 = s3
	server.emailService = email
	server.storage = authStore
	server.redis = redis
	return server
}
