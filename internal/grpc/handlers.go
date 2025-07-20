package grpc

import (
	"context"
	"log"

	"github.com/MagicRodri/grpc_with_go/pkg/helloworld"
	"github.com/MagicRodri/grpc_with_go/pkg/status"
)

func (s *Server) SayHello(_ context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *Server) SetStatus(_ context.Context, in *status.StatusMessage) (*status.StatusResponse, error) {
	log.Printf("Received: %v", in.GetUuid())
	return &status.StatusResponse{Uuid: in.GetUuid(), Message: "Status set", Code: 0}, nil
}

func (s *Server) GetStatus(_ context.Context, in *status.StatusRequest) (*status.StatusResponse, error) {
	log.Printf("Received: %v", in.GetUuid())
	return &status.StatusResponse{Uuid: in.GetUuid(), Message: "Status retrieved", Code: 0}, nil
}
func (s *Server) DeleteStatus(_ context.Context, in *status.StatusRequest) (*status.StatusResponse, error) {
	log.Printf("Received: %v", in.GetUuid())
	return &status.StatusResponse{Uuid: in.GetUuid(), Message: "Status deleted", Code: 0}, nil
}
