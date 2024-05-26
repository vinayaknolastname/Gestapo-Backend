package service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
)

// merchantService serves gRPC requests for our e-commerce service.
type merchantService struct {
	proto.UnimplementedMerchantServiceServer
	config  *config.Config
	log     logger.Logger
	s3      *s3.S3Service
	storage *db.MerchantStore
	token   token.Maker
}

// NewMerchantService creates a new gRPC server.
func NewMerchantService(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *merchantService {
	server := &merchantService{
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

	merchantStore := db.NewMerchantStore(storage)

	server.s3 = s3
	server.storage = merchantStore
	return server
}
