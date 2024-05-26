package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/api/proto"
	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/service"
	"github.com/akmal4410/gestapo/pkg/helpers/interceptor"
	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCService(ctx context.Context, storage *database.Storage, config *config.Config, log logger.Logger) error {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.LogFatal("Error while Initializing NewJWTMaker %w", err)
	}
	service := service.NewMerchantService(storage, config, log, tokenMaker)
	interceptor := interceptor.NewInterceptor(tokenMaker, log)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.AccessMiddleware(),
			interceptor.MerchantRoleMiddleware(),
		),
	)

	proto.RegisterMerchantServiceServer(grpcServer, service)
	log.LogInfo("Registreing for reflection")
	reflection.Register(grpcServer)
	port := ":" + config.ServerAddress.Merchant.Port

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.LogError("error in listening to port", port, "error:", err)
		return err
	}
	//graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range c {
			log.LogInfo("shutting down grpc server....")
			grpcServer.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.LogInfo("Start gRPC server at ", lis.Addr().String())
	return grpcServer.Serve(lis)
}
