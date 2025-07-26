package server

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CombinedGRPCServer struct {
	server *grpc.Server
	logger *zap.Logger
}

func NewCombinedGRPCServer(logger *zap.Logger) *CombinedGRPCServer {
	s := &CombinedGRPCServer{
		server: grpc.NewServer(),
		logger: logger.Named("combined-grpc-server"),
	}

	// TODO: Register services once they're fixed
	// pb.RegisterPropertyServiceServer(s.server, s.propertyServer)

	return s
}

func (s *CombinedGRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	s.logger.Info("gRPC server starting", zap.Int("port", port))
	return s.server.Serve(lis)
}

func (s *CombinedGRPCServer) Stop() {
	s.logger.Info("stopping gRPC server")
	s.server.GracefulStop()
}
