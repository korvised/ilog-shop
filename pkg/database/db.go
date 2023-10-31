package database

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func DbConn(c context.Context, cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.Url))
	if err != nil {
		log.Fatalf("Error: connecting to mongodb: %s \n", err.Error())
	}

	// Ping MongoDB
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Error: pinging mongodb: %s \n", err.Error())
	}

	return client
}
