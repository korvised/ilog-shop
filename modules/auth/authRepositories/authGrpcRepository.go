package authRepositories

import (
	"context"
	"errors"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func (r *authRepository) FindCredential(c context.Context, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(r.cfg.Grpc.PlayerUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %v \n", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)
	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %v \n", err)
		return nil, errors.New("error: invalid credential")
	}

	return result, nil
}

func (r *authRepository) FindOnePlayerProfileToRefresh(c context.Context, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(r.cfg.Grpc.PlayerUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %v \n", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)
	result, err := conn.Player().FindOnePlayerProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh failed: %v \n", err)
		return nil, errors.New("error: player not found")
	}

	return result, nil
}

func (r *authRepository) FindRoleCount(c context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("roles")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: FindRoleCount: %v \n", err)
		return -1, errors.New("error: find role count failed")
	}

	return count, nil
}
