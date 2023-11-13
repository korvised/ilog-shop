package grpcconn

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func (g *grpcAuth) unaryAuthorization(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("Error: gRPC Metadata not found")
		return nil, errors.New("error: metadata not found")
	}

	authHeader, ok := md["auth"]
	if !ok {
		log.Println("Error: gRPC Authorization header not found")
		return nil, errors.New("error: authorization header not found")
	}

	if len(authHeader) == 0 {
		log.Println("Error: gRPC Authorization header invalid")
		return nil, errors.New("error: authorization header invalid")
	}

	claims, err := jwtauth.ParseToken(g.secretKey, authHeader[0])
	if err != nil {
		log.Printf("Error: gRPC failed to parse token: %v \n", err)
		return nil, errors.New("error: failed to parse token")
	}
	log.Printf("Claims: %v \n", claims)

	return handler(ctx, req)
}
