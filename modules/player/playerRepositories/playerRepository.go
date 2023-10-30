package playerRepositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PlayerRepositoryService interface {
	}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayerRepository(db *mongo.Client) PlayerRepositoryService {
	return &playerRepository{db}
}

func (r *playerRepository) playerDbConn(c context.Context) *mongo.Database {
	return r.db.Database("player_db")
}
