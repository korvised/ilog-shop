package migration

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func inventoryDbConn(ctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(ctx, cfg).Database("inventory_db")
}

func InventoryMigrate(ctx context.Context, cfg *config.Config) {
	db := inventoryDbConn(ctx, cfg)
	defer db.Client().Disconnect(ctx)

	col := db.Collection("players_inventory")

	indexs, _ := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"player_id", 1}, {"item_id", 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	col = db.Collection("players_inventory_queue")

	results, err := col.InsertOne(ctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate inventory completed: ", results)
}
