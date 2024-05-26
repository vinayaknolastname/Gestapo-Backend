package service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	s3 "github.com/akmal4410/gestapo/pkg/service/s3_service"
)

type userService struct {
	*proto.UnimplementedUserServieServer
	config  *config.Config
	log     logger.Logger
	s3      *s3.S3Service
	storage *db.UserStore
	token   token.Maker
}

// NewUserService creates a new gRPC server.
func NewUserService(storage *database.Storage, config *config.Config, log logger.Logger, tokenMaker token.Maker) *userService {
	server := &userService{
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
	userStore := db.NewUserStore(storage)
	server.s3 = s3
	server.storage = userStore
	return server
}
