package middlewareRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
	"time"
)

type (
	MiddlewareRepositoryService interface {
		FindOneCredential(c context.Context, accessToken string) error
		FineRoleCount(c context.Context) (int64, error)
	}

	middlewareRepository struct {
		cfg *config.Config
	}
)

func NewMiddlewareRepository(cfg *config.Config) MiddlewareRepositoryService {
	return &middlewareRepository{cfg: cfg}
}

func (r *middlewareRepository) FindOneCredential(c context.Context, accessToken string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(r.cfg.Grpc.AuthUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %v \n", err)
		return errors.New("error: gRPC client connection failed")
	}

	result, err := conn.Auth().GetCredential(ctx, &authPb.CredentialReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("Error: GetCredential: %v \n", err)
		return errors.New("error: invalid credential")
	}

	if !result.GetIsValid() {
		log.Println("Error: invalid credential")
		return errors.New("error: invalid credential")
	}

	return nil
}

func (r *middlewareRepository) FineRoleCount(c context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(r.cfg.Grpc.AuthUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %v \n", err)
		return -1, errors.New("error: gRPC client connection failed")
	}

	result, err := conn.Auth().GetRolesCount(ctx, &authPb.RolesCountReq{})
	if err != nil {
		log.Printf("Error: GetCredential: %v \n", err)
		return -1, errors.New("error: get role count failed")
	}

	return result.GetCount(), nil
}
