package authRepositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepositoryService interface {
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db}
}

func (r *authRepository) authDbConn(c context.Context) *mongo.Database {
	return r.db.Database("auth-db")
}
