package service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
)

// productService serves gRPC requests for our e-commerce service.
type productService struct {
	proto.UnimplementedProductServiceServer
	config  *config.Config
	log     logger.Logger
	s3      *s3.S3Service
	storage *db.ProductStore
	token   token.Maker
}

// NewProductService creates a new gRPC server.
func NewProductService(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *productService {
	server := &productService{
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

	productStore := db.NewProductStore(storage)

	server.s3 = s3
	server.storage = productStore
	return server
}
