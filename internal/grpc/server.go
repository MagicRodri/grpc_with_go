package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MagicRodri/grpc_with_go/config"
	"github.com/MagicRodri/grpc_with_go/pkg/generated/helloworld"
	"github.com/MagicRodri/grpc_with_go/pkg/generated/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents a gRPC server.
type Server struct {
	helloworld.GreeterServer
	status.StatusServiceServer
	grpcServer *grpc.Server
	address    string
}

// NewServer creates a new gRPC server instance.
func NewServer(cfg *config.GrpcConfig) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		address:    cfg.Address,
	}
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}
	helloworld.RegisterGreeterServer(s.grpcServer, s)
	status.RegisterStatusServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)
	fmt.Printf("gRPC server listening on %s\n", s.address)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	fmt.Println("gRPC server stopped")
}

// Serve starts the gRPC server and blocks until it is stopped.
func (s *Server) Serve(ctx context.Context) error {
	err := s.Start()
	if err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	// Wait for the context to be done before stopping the server
	<-ctx.Done()
	s.Stop()
	return nil
}
