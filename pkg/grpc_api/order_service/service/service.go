package service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/order_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
)

// orderService serves gRPC requests for our e-commerce service.
type orderService struct {
	proto.UnimplementedOrderServiceServer
	config  *config.Config
	log     logger.Logger
	s3      *s3.S3Service
	storage *db.OrderStore
	token   token.Maker
}

// NewOrderService creates a new gRPC server.
func NewOrderService(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *orderService {
	server := &orderService{
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

	orderStore := db.NewOrderStore(storage)

	server.s3 = s3
	server.storage = orderStore
	return server
}
