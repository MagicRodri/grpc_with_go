package main

import (
	"context"
	"flag"
	"log"

	"github.com/MagicRodri/grpc_with_go/config"
	"github.com/MagicRodri/grpc_with_go/internal/grpc"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to the configuration file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Error loading config: %v\n", err)
		return
	}
	server := grpc.NewServer(&cfg.GRPC)
	if err := server.Serve(context.Background()); err != nil {
		log.Printf("Error starting gRPC server: %v\n", err)
		return
	}
	defer server.Stop()
	log.Println("gRPC server started successfully")
}
