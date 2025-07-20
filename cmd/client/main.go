// client main package
package main

import (
	"context"
	"fmt"

	"github.com/MagicRodri/grpc_with_go/config"
	"github.com/MagicRodri/grpc_with_go/internal/grpc"
)

func main() {
	configPath := "config/config.yaml"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	client, err := grpc.NewClient(cfg.GRPC.Address)
	if err != nil {
		fmt.Printf("Error creating gRPC client: %v\n", err)
		return
	}
	defer client.Close()
	response, err := client.SayHello(context.Background(), "World")
	if err != nil {
		fmt.Printf("Error saying hello: %v\n", err)
		return
	}
	fmt.Printf("Response from server: %v\n", response)
}
