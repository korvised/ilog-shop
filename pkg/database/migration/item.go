package migration

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/item"
	"github.com/korvised/ilog-shop/pkg/database"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func itemDbConn(ctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(ctx, cfg).Database("item_db")
}

func ItemMigrate(ctx context.Context, cfg *config.Config) {
	db := itemDbConn(ctx, cfg)
	defer db.Client().Disconnect(ctx)

	col := db.Collection("items")
	indexs, _ := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"title", 1}}},
	})
	log.Println(indexs)

	documents := func() []any {
		roles := []*item.Item{
			{
				Title:       "Diamond Sword",
				Price:       1000,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      100,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
			{
				Title:       "Iron Sword",
				Price:       500,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      50,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
			{
				Title:       "Wooden Sword",
				Price:       100,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      20,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(ctx, documents, nil)
	if err != nil {
		log.Fatalf("Error migrate item: %v", err)
	}
	log.Println("Migrate item completed: ", results)
}
