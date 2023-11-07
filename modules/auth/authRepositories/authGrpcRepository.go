package authRepositories

import (
	"context"
	"errors"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
	"time"
)

func (r *authRepository) CredentialSearch(c context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %v \n", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %v \n", err)
		return nil, errors.New("error: invalid credential")
	}

	return result, nil
}

func (r *authRepository) FindOnePlayerProfileToRefresh(c context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %v \n", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	result, err := conn.Player().FindOnePlayerProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh failed: %v \n", err)
		return nil, errors.New("error: player not found")
	}

	return result, nil
}
