package paymentRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/inventory"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/models"
	"github.com/korvised/ilog-shop/modules/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type (
	PaymentRepositoryService interface {
		FindItemInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
		DockedPlayerMoney(c context.Context, req *player.CreatePlayerTransactionReq) error
		RollbackTransaction(c context.Context, req *player.RollbackPlayerTransactionReq) error
		AddPlayItem(c context.Context, req *inventory.UpdateInventoryReq) error
		RollbackAddPlayItem(c context.Context, req *inventory.RollbackPlayerInventoryReq) error
	}

	paymentRepository struct {
		db  *mongo.Client
		cfg *config.Config
	}
)

func NewPaymentRepository(db *mongo.Client, cfg *config.Config) PaymentRepositoryService {
	return &paymentRepository{db, cfg}
}

func (r *paymentRepository) paymentDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}

func (r *paymentRepository) FindOffset(c context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: FindOffset: %v \n", err)
		return -1, errors.New("error: payment offset not found")
	}

	return result.Offset, nil
}

func (r *paymentRepository) UpsertOffset(c context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	_, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpsertOffset: %v \n", err)
		return errors.New("error: upsert payment offset failed")
	}

	return nil
}
