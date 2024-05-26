package user_service

import (
	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/user_service/protocol/grpc"
	"github.com/akmal4410/gestapo/pkg/helpers/service_helper"
)

const (
	serviceName = "User Service"
	logFileName = "user_service"
)

func RunServer() error {
	ctx, log := service_helper.InitializeService(serviceName, logFileName)
	config, err := config.LoadConfig("configs")
	if err != nil {
		log.LogFatal("Cannot load configuration:", err)
	}
	log.LogInfo("Config file loaded.")
	store, err := database.NewStorage(config.Database)
	if err != nil {
		log.LogFatal("Cannot connect to Database", err)
	}
	log.LogInfo("Database connection successful")

	err = grpc.RunGRPCService(ctx, store, &config, log)
	if err != nil {
		log.LogFatal("Cannot start server :", err)
		return err
	}

	select {
	case <-ctx.Done():
		log.LogError(ctx.Err())
		break
	}
	log.LogError(serviceName, "shutdown")
	return nil
}
