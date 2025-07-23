package client

import (
	"context"
	"fmt"

	"time"

	status_service "github.com/MagicRodri/grpc_with_go/pkg/generated/status"
	"github.com/MagicRodri/grpc_with_go/pkg/logger"
	"github.com/MagicRodri/grpc_with_go/pkg/manager"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatusClient struct {
	conn   *grpc.ClientConn
	client status_service.StatusServiceClient
	log    logger.LoggerInterface
	cfg    *StatusServiceConfig
}

type Status struct {
	Uuid      string
	Timestamp time.Time
}

func NewStatusClient(log logger.LoggerInterface, cfg *StatusServiceConfig) *StatusClient {
	return &StatusClient{
		log: log,
		cfg: cfg,
	}
}

func (sc *StatusClient) Initialize(conn *grpc.ClientConn) error {
	sc.conn = conn
	sc.client = status_service.NewStatusServiceClient(conn)
	sc.log.Info("Statistics gRPC client initialized successfully")
	return nil
}

func (sc *StatusClient) Close() error {
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}

func (sc *StatusClient) GetName() string {
	return sc.cfg.Name
}

func (sc *StatusClient) GetHost() string {
	return sc.cfg.Host
}

func (sc *StatusClient) SendStatusMessage(status *Status) error {
	if sc.client == nil {
		return fmt.Errorf("statistics gRPC client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sc.cfg.Timeout)*time.Second)
	defer cancel()

	req := &status_service.StatusMessage{
		Uuid:      status.Uuid,
		Timestamp: timestamppb.New(status.Timestamp),
	}

	res, err := sc.client.SetStatus(ctx, req)
	if err != nil {
		sc.log.Error("failed to record status: %v", err)
		return fmt.Errorf("failed to record status: %w", err)
	}

	sc.log.Debug("Successfully recorded status: %+v, response: %+v", status, res)
	return nil
}

// helper function to get the StatusClient from the global client manager
func GetStatusClient(name string) (*StatusClient, error) {
	if manager.GlobalGrpcClientManager == nil {
		return nil, fmt.Errorf("global client manager is not initialized")
	}

	client, err := manager.GlobalGrpcClientManager.GetClient(name)
	if err != nil {
		return nil, err
	}

	statusClient, ok := client.(*StatusClient)
	if !ok {
		return nil, fmt.Errorf("client is not a status client")
	}

	return statusClient, nil
}
