package inventoryRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/inventory"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/models"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type (
	InventoryRepositoryService interface {
		InsertOnePlayerItem(c context.Context, req *inventory.Inventory) (primitive.ObjectID, error)
		FindItemInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindPlayItems(c context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error)
		CountPlayerItems(c context.Context, filter primitive.D) (int64, error)
		FindOnePlayerItem(c context.Context, PlayerID, ItemID string) error
		FindOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
		AddPlayerItemRes(_ context.Context, req *payment.PaymentTransferRes) error
		RemovePlayerItemRes(_ context.Context, req *payment.PaymentTransferRes) error
		DeleteOneInventory(c context.Context, inventoryID string) error
		DeleteOnePlayerItem(c context.Context, PlayerID, ItemID string) error
	}

	inventoryRepository struct {
		db  *mongo.Client
		cfg *config.Config
	}
)

func NewInventoryRepository(db *mongo.Client, cfg *config.Config) InventoryRepositoryService {
	return &inventoryRepository{db, cfg}
}

func (r *inventoryRepository) inventoryDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}

func (r *inventoryRepository) InsertOnePlayerItem(c context.Context, req *inventory.Inventory) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerItem: %s\n", err.Error())
		return primitive.NilObjectID, errors.New("error: insert player item failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *inventoryRepository) FindPlayItems(c context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindPlayItems: %v \n", err)
		return nil, errors.New("error: player item not found")
	}

	results := make([]*inventory.Inventory, 0)
	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err = cursors.Decode(&result); err != nil {
			log.Printf("Error: FindPlayItems: %v \n", err)
			return nil, errors.New("error: player item not found")
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountPlayerItems(c context.Context, filter primitive.D) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: CountItems: %s\n", err.Error())
		return -1, errors.New("error: count items failed")
	}

	return count, nil
}

func (r *inventoryRepository) FindOnePlayerItem(c context.Context, PlayerID, ItemID string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result := new(inventory.Inventory)
	if err := col.FindOne(ctx, bson.M{"player_id": PlayerID, "item_id": ItemID}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerItem: %v \n", err)
		return errors.New("error: player item not found")
	}

	return nil
}

func (r *inventoryRepository) FindOffset(c context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: FindOffset: %v \n", err)
		return -1, errors.New("error: inventory offset not found")
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpsertOffset(c context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	_, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpsertOffset: %v \n", err)
		return errors.New("error: upsert inventory offset failed")
	}

	return nil
}

func (r *inventoryRepository) DeleteOneInventory(c context.Context, inventoryID string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	if _, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(inventoryID)}); err != nil {
		log.Printf("Error: DeleteOnePlayerItem: %s\n", err.Error())
		return errors.New("error: delete player item failed")
	}

	return nil
}

func (r *inventoryRepository) DeleteOnePlayerItem(c context.Context, PlayerID, ItemID string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	if _, err := col.DeleteOne(ctx, bson.M{"player_id": PlayerID, "item_id": ItemID}); err != nil {
		log.Printf("Error: DeleteOnePlayerItem: %s\n", err.Error())
		return errors.New("error: delete player item failed")
	}

	return nil
}
