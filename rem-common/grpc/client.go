package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Dial establece la conexión gRPC con keep-alive sensato.
// – ctx con timeout de 5 s  ➜ aborta si no conecta.
// – usa credenciales inseguras *solo* para dev/local.
func Dial(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: TLS en prod
		grpc.WithBlock(), // bloqueo hasta que conecte o ctx venza
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
	)
}
