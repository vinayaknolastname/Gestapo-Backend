package service

import (
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/admin_service/db"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
)

type adminService struct {
	proto.UnimplementedAdminServiceServer
	storage *db.AdminStore
	log     logger.Logger
}

// NewAuthenticationService creates a new gRPC server.
func NewAdminService(storage *database.Storage, log logger.Logger) *adminService {
	server := &adminService{
		log: log,
	}

	authStore := db.NewAdminStore(storage)

	server.storage = authStore
	return server
}
