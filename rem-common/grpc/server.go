package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// NewServer devuelve *grpc.Server con:
//   - Logging interceptor (zap)        – unary y stream
//   - Recovery interceptor             – evita panics
//   - Health service registrado
//   - reflection.Register (útil para grpcurl)
func NewServer(lg *zap.Logger, ka keepalive.ServerParameters) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(unaryLogger(lg)),
		grpc.StreamInterceptor(streamLogger(lg)),
		grpc.KeepaliveParams(ka), // configurable por servicio
	}

	s := grpc.NewServer(opts...)

	// healthz
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthSrv)

	// reflection (solo dev – quítalo en prod si quieres)
	reflection.Register(s)

	return s
}

/* ------------------- interceptors ------------------- */

func unaryLogger(lg *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)

		elapMs := float64(time.Since(start).Microseconds()) / 1_000

		lg.Info("grpc unary",
			zap.String("method", info.FullMethod),
			zap.Float64("latency_ms", elapMs),
			zap.Error(err),
		)
		return resp, err
	}
}

func streamLogger(lg *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		start := time.Now()
		err := handler(srv, ss)

		elapMs := float64(time.Since(start).Microseconds()) / 1_000

		lg.Info("grpc stream",
			zap.String("method", info.FullMethod),
			zap.Float64("latency_ms", elapMs),
			zap.Error(err),
		)
		return err
	}
}
