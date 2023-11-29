package playerRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/models"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type (
	PlayerRepositoryService interface {
		FindOnePlayerProfile(c context.Context, playerID string) (*player.PlayerProfileBson, error)
		FindOnePlayerSavingAccount(c context.Context, playerID string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(c context.Context, email string) (*player.Player, error)
		FindOnePlayerProfileToRefresh(c context.Context, playerID string) (*player.Player, error)
		IsUniquePlayer(c context.Context, email, username string) bool
		InsertOnePlayer(c context.Context, req *player.Player) (primitive.ObjectID, error)
		InsertOnePlayerTransaction(c context.Context, req *player.PlayerTransaction) (primitive.ObjectID, error)
		FindOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
		DockedPlayerMoneyRes(_ context.Context, req *payment.PaymentTransferRes) error
		DeleteOnePlayerTransaction(c context.Context, transactionID string) error
	}

	playerRepository struct {
		db  *mongo.Client
		cfg *config.Config
	}
)

func NewPlayerRepository(db *mongo.Client, cfg *config.Config) PlayerRepositoryService {
	return &playerRepository{db, cfg}
}

func (r *playerRepository) playerDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerRepository) IsUniquePlayer(c context.Context, email, username string) bool {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.Player)
	if err := col.
		FindOne(ctx, bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}).
		Decode(result); err != nil {
		return true
	}

	return false
}

func (r *playerRepository) InsertOnePlayer(c context.Context, req *player.Player) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	playerID, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: insert player failed: %v", err)
		return primitive.NilObjectID, errors.New("error: insert player failed")
	}

	return playerID.InsertedID.(primitive.ObjectID), nil
}

func (r *playerRepository) FindOnePlayerProfile(c context.Context, playerID string) (*player.PlayerProfileBson, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.PlayerProfileBson)
	if err := col.
		FindOne(
			ctx,
			bson.M{"_id": utils.ConvertToObjectId(playerID)},
			options.FindOne().SetProjection(bson.M{
				"_id":        1,
				"email":      1,
				"username":   1,
				"created_at": 1,
				"updated_at": 1,
			}),
		).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerProfile: %v", result)
		return nil, errors.New("error: player profile not found")
	}

	return result, nil
}

func (r *playerRepository) InsertOnePlayerTransaction(c context.Context, req *player.PlayerTransaction) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerTransaction: %v", err)
		return primitive.NilObjectID, errors.New("error: insert player transactions failed")
	}
	log.Printf("Result: InsertOnePlayerTransaction: %v", result.InsertedID)

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *playerRepository) FindOnePlayerSavingAccount(c context.Context, playerID string) (*player.PlayerSavingAccount, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	filter := bson.A{
		bson.D{{"$match", bson.D{{"player_id", playerID}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$player_id"},
					{"balance", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"player_id", "$_id"},
					{"balance", bson.D{{"$toDouble", "$balance"}}},
				},
			},
		},
	}

	cursors, err := col.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Error: FindOnePlayerSavingAccount: %v \n", err)
		return nil, errors.New("error: find player saving account failed")
	}

	result := new(player.PlayerSavingAccount)
	for cursors.Next(ctx) {
		if err = cursors.Decode(result); err != nil {
			log.Printf("Error: Decode FindOnePlayerSavingAccount: %v \n", err)
			return nil, errors.New("error: find player saving account failed")
		}
	}

	return result, nil
}

func (r *playerRepository) FindOnePlayerCredential(c context.Context, email string) (*player.Player, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	document := new(player.Player)

	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(document); err != nil {
		log.Printf("Error: FindOnePlayerCredential failed: %v \n", err)
		return nil, errors.New("error: player not found")
	}

	return document, nil
}

func (r *playerRepository) FindOnePlayerProfileToRefresh(c context.Context, playerID string) (*player.Player, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	document := new(player.Player)
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(playerID)}).Decode(document); err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh failed: %v \n", err)
		return nil, errors.New("error: player not found")
	}

	return document, nil
}

func (r *playerRepository) FindOffset(c context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: FindOffset: %v \n", err)
		return -1, errors.New("error: player offset not found")
	}

	return result.Offset, nil
}

func (r *playerRepository) UpsertOffset(c context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions_queue")

	_, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpsertOffset: %v \n", err)
		return errors.New("error: upsert player offset failed")
	}

	return nil
}

func (r *playerRepository) DeleteOnePlayerTransaction(c context.Context, transactionID string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	if _, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(transactionID)}); err != nil {
		log.Printf("Error: DeleteOnePlayerTransaction: %s\n", err.Error())
		return errors.New("error: delete player transaction failed")
	}

	log.Printf("Delete player transaction %s success", transactionID)

	return nil
}
