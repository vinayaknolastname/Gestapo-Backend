package server

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
)

type RestServer struct {
	log     logger.Logger
	s3      *s3.S3Service
	storage *db.MerchantStore
	token   token.Maker
}

// NewRestServer creates a new server for handling http request.
func NewRestServer(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *RestServer {
	server := &RestServer{
		log:   log,
		token: tokenMaker,
	}
	s3 := s3.NewS3Service(
		config.AwsS3.BucketName,
		config.AwsS3.Region,
		config.AwsS3.AccessKey,
		config.AwsS3.SecretKey,
	)

	merchantStore := db.NewMerchantStore(storage)

	server.s3 = s3
	server.storage = merchantStore
	return server
}
