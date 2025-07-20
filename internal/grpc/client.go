// Package grpc provides a gRPC client for the Greeter service.
package grpc

import (
	"context"
	"fmt"

	"github.com/MagicRodri/grpc_with_go/pkg/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a gRPC client for the Greeter service.
type Client struct {
	conn   *grpc.ClientConn
	client helloworld.GreeterClient
}

// NewClient creates a new gRPC client for the Greeter service.
func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	client := helloworld.NewGreeterClient(conn)
	return &Client{conn: conn, client: client}, nil
}

// SayHello sends a hello request to the gRPC server.
func (c *Client) SayHello(ctx context.Context, name string) (string, error) {
	req := &helloworld.HelloRequest{Name: name}
	res, err := c.client.SayHello(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to say hello: %w", err)
	}
	return res.Message, nil
}

// Close closes the gRPC client connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
