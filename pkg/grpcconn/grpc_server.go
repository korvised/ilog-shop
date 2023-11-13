package grpcconn

import (
	"github.com/korvised/ilog-shop/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{secretKey: cfg.ApiSecretKey}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: failed to listen: %v", err)
	}

	return grpcServer, lis
}
