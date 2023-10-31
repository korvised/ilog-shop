package migration

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth"
	"github.com/korvised/ilog-shop/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func authDbConn(ctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(ctx, cfg).Database("auth_db")
}

func AuthMigrate(ctx context.Context, cfg *config.Config) {
	db := authDbConn(ctx, cfg)
	defer db.Client().Disconnect(ctx)

	// auth
	col := db.Collection("auth")

	indexs, _ := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
		{Keys: bson.D{{"refresh_token", 1}}},
	})
	log.Println(indexs)

	// roles
	col = db.Collection("roles")

	indexs, _ = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"code", 1}}},
	})
	log.Println(indexs)

	// role data
	documents := func() []any {
		roles := []*auth.Role{
			{
				Title: "user",
				Code:  0,
			},
			{
				Title: "admin",
				Code:  1,
			},
		}

		docs := make([]any, 0)
		for _, role := range roles {
			docs = append(docs, role)
		}
		return docs
	}()

	result, err := col.InsertMany(ctx, documents)
	if err != nil {
		log.Fatalf("Error migrate auth: %v", err)
	}
	log.Println("Migrate auth completed: ", result.InsertedIDs)
}
