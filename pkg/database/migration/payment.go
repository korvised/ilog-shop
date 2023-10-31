package migration

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func paymentDbConn(ctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(ctx, cfg).Database("payment_db")
}

func PaymentMigrate(ctx context.Context, cfg *config.Config) {
	db := paymentDbConn(ctx, cfg)
	defer db.Client().Disconnect(ctx)

	col := db.Collection("payment_queue")

	results, err := col.InsertOne(ctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate payment completed: ", results)
}
