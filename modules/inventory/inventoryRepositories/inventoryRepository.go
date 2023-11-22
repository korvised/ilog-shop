package inventoryRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/inventory"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type (
	InventoryRepositoryService interface {
		FindItemInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindPlayItems(c context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error)
		CountPlayerItems(c context.Context, filter primitive.D) (int64, error)
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
