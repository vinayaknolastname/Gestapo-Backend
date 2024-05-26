package main

import (
	"fmt"
	"os"

	"github.com/akmal4410/gestapo/pkg/grpc_api/grpc_gateway"
)

func main() {
	err := grpc_gateway.RunGateway()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
