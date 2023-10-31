package migration

import (
	"context"
	"fmt"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/database"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func playerDbConn(ctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(ctx, cfg).Database("player_db")
}

func PlayerMigrate(ctx context.Context, cfg *config.Config) {
	fmt.Println("---- Migrate player -----")

	db := playerDbConn(ctx, cfg)
	defer db.Client().Disconnect(ctx)

	// player transactions
	col := db.Collection("player_transactions")
	indexs, _ := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
	})
	log.Println(indexs)

	// player
	col = db.Collection("players")
	indexs, _ = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"email", 1}}},
	})
	log.Println(indexs)

	// Seed data
	documents := func() []any {
		roles := []*player.Player{
			{
				Email: "player001@sekai.com",
				Password: func() string {
					// Hashing password
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(hashedPassword)
				}(),
				Username: "Player001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email: "player002@sekai.com",
				Password: func() string {
					// Hashing password
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(hashedPassword)
				}(),
				Username: "Player002",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email: "player003@sekai.com",
				Password: func() string {
					// Hashing password
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(hashedPassword)
				}(),
				Username: "Player003",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email: "admin001@sekai.com",
				Password: func() string {
					// Hashing password
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(hashedPassword)
				}(),
				Username: "Player003",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(ctx, documents)
	if err != nil {
		log.Fatalf("Error migrate player: %v", err)
	}
	log.Println("Migrate player completed: ", results.InsertedIDs)

	playerTransactions := make([]any, 0)
	for _, p := range results.InsertedIDs {
		playerTransactions = append(playerTransactions, &player.PlayerTransaction{
			PlayerID:  "player:" + p.(primitive.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.LocalTime(),
		})
	}

	col = db.Collection("player_transactions")
	results, err = col.InsertMany(ctx, playerTransactions, nil)
	if err != nil {
		log.Fatalf("Error migrate player_transactions: %v", err)
	}
	log.Println("Migrate player_transactions completed: ", results.InsertedIDs)

	col = db.Collection("player_transactions_queue")
	result, err := col.InsertOne(ctx, bson.M{"offset": -1}, nil)
	if err != nil {
		log.Fatalf("Error migrate player_transactions_queue: %v", err)
	}
	log.Println("Migrate player_transactions_queue completed: ", result)
}
