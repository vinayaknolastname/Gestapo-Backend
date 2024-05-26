package service_helper

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/akmal4410/gestapo/pkg/helpers/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitializeService func is called to set the log service and also listen to interrupt signals
func InitializeService(serviceName string, logFileName string) (context.Context, logger.Logger) {
	//listen to interrupts or process termination
	//first we create a channel to listen to interrupt's or kill signals
	c := make(chan os.Signal, 1)
	//we pass the channel to the single notify func
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	//we create a context that we will use the program
	ctx, cancel := context.WithCancel(context.Background())
	//we listen to the channel and trigger cancel func
	go func() {
		<-c
		cancel()
	}()
	log := logger.NewLogrusLogger(logFileName)
	log.LogInfo(serviceName, "has started")
	return ctx, log
}

func ConnectEndpoints(address, serviceName string, log logger.Logger) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.LogError("connection to", serviceName, "(", address, ") failed. Error details:", err)
		return nil, err
	}
	return conn, nil
}
