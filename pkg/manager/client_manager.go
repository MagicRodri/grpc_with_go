package manager

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/MagicRodri/grpc_with_go/pkg/logger"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxRetry     uint          = 5
	retryTimeOut time.Duration = time.Millisecond * 100
)

type GrpcClientInterface interface {
	Initialize(conn *grpc.ClientConn) error
	Close() error
	GetName() string
	GetHost() string
}

// GrpcClientManager manages gRPC clients
// It allows registering, retrieving, and closing clients.
type GrpcClientManager struct {
	clients map[string]GrpcClientInterface
	mutex   sync.RWMutex
	log     logger.LoggerInterface
}

func NewGrpcClientManager(log logger.LoggerInterface) *GrpcClientManager {
	return &GrpcClientManager{
		clients: make(map[string]GrpcClientInterface),
		log:     log,
	}
}

func (cm *GrpcClientManager) RegisterClient(client GrpcClientInterface) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.clients[client.GetName()]; exists {
		return fmt.Errorf("client '%s' is already registered", client.GetName())
	}

	conn, err := cm.createConnection(client.GetHost())
	if err != nil {
		return fmt.Errorf("failed to create connection for client '%s': %w", client.GetName(), err)
	}

	if err := client.Initialize(conn); err != nil {
		conn.Close()
		return fmt.Errorf("failed to initialize client '%s': %w", client.GetName(), err)
	}

	cm.clients[client.GetName()] = client
	cm.log.Info("Successfully registered gRPC client: %s", client.GetName())
	return nil
}

func (cm *GrpcClientManager) GetClient(name string) (GrpcClientInterface, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	client, exists := cm.clients[name]
	if !exists {
		return nil, fmt.Errorf("client '%s' not found", name)
	}
	return client, nil
}

func (cm *GrpcClientManager) CloseAll() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for name, client := range cm.clients {
		if err := client.Close(); err != nil {
			cm.log.Error("Error closing client '%s': %v", name, err)
		} else {
			cm.log.Info("Closed gRPC client: %s", name)
		}
	}
	cm.clients = make(map[string]GrpcClientInterface)
}

func (cm *GrpcClientManager) ListClients() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	names := make([]string, 0, len(cm.clients))
	for name := range cm.clients {
		names = append(names, name)
	}
	return names
}

func (cm *GrpcClientManager) createConnection(address string) (*grpc.ClientConn, error) {
	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		retry.WithMax(maxRetry),
		retry.WithBackoff(retry.BackoffLinear(retryTimeOut)),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(cm.logInterceptor(), logOpts...),
			retry.UnaryClientInterceptor(retryOpts...),
		),
	}

	conn, err := grpc.NewClient(address, grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	return conn, nil
}

func (cm *GrpcClientManager) logInterceptor() grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		cm.log.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// GlobalGrpcClientManager - a singleton instance of GrpcClientManager
var GlobalGrpcClientManager *GrpcClientManager

func InitGlobalGrpcClientManager(log logger.LoggerInterface) {
	if GlobalGrpcClientManager == nil {
		GlobalGrpcClientManager = NewGrpcClientManager(log)
	}
}
