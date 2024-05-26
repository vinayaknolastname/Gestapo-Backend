package grpc_gateway

import (
	"context"
	"net/http"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/grpc_gateway/server"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/gorilla/handlers"
)

const (
	serviceName = "gRPC Gateway"
	logFileName = "grpc_gateway"
)

func RunGateway() error {
	log := logger.NewLogrusLogger(logFileName)
	log.LogInfo(serviceName, "has started")

	config, err := config.LoadConfig("configs")
	if err != nil {
		log.LogFatal("Cannot load configuration:", err)
	}
	log.LogInfo("Config file loaded.")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gMux, err := newGateway(ctx, log, config)
	if err != nil {
		log.LogError("error in newGateway :", err)
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", gMux))

	//-----------------ONlY FOR REST API(Image handling)--------------------------
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.LogFatal("Error while Initializing NewJWTMaker %w", err)
	}
	store, err := database.NewStorage(config.Database)
	if err != nil {
		log.LogFatal("Cannot connect to Database", err)
	}
	log.LogInfo("Database connection successful")
	server := server.NewRestServer(store, &config, log, tokenMaker)
	server.SetupRouter(mux)
	//------------------------------------------------------------------------------

	return http.ListenAndServe(":"+config.ServerAddress.Gateway,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "User-Agent"}),
			handlers.ExposedHeaders([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(mux),
	)
}
