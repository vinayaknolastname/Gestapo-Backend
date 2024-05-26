package main

import (
	"fmt"
	"os"

	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service"
)

func main() {
	err := merchant_service.RunServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
