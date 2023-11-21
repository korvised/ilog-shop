package itemRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/modules/item"
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
	ItemRepositoryService interface {
		IsUniqueItem(c context.Context, title string) bool
		InsertOneItem(c context.Context, req *item.Item) (primitive.ObjectID, error)
		FindOneItem(c context.Context, itemID string) (*item.Item, error)
		FindManyItems(c context.Context, filter primitive.D, opts []*options.FindOptions) ([]*item.ItemShowCase, error)
		CountItems(c context.Context, filter primitive.D) (int64, error)
		UpdateItem(c context.Context, itemID string, req primitive.M) error
		UpdateItemStatus(c context.Context, itemID string, isActive bool) error
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{db}
}

func (r *itemRepository) itemDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) IsUniqueItem(c context.Context, title string) bool {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(player.Player)
	if err := col.FindOne(ctx, bson.M{"title": title}).Decode(result); err != nil {
		log.Printf("Error: IsUniqueItem: %s\n", err.Error())
		return true
	}

	return false
}

func (r *itemRepository) InsertOneItem(c context.Context, req *item.Item) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	itemID, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneItem: %s\n", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one item failed")
	}

	return itemID.InsertedID.(primitive.ObjectID), nil
}

func (r *itemRepository) FindOneItem(c context.Context, itemID string) (*item.Item, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemID)}).Decode(result); err != nil {
		log.Printf("Error: FindOneItem: %s\n", err.Error())
		return nil, errors.New("error: find one item failed")
	}

	return result, nil
}

func (r *itemRepository) FindManyItems(c context.Context, filter primitive.D, opts []*options.FindOptions) ([]*item.ItemShowCase, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindManyItems: %s\n", err.Error())
		return nil, errors.New("error: find many items failed")
	}
	results := make([]*item.ItemShowCase, 0)
	for cursors.Next(ctx) {
		result := new(item.Item)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindManyItems: %s\n", err.Error())
			return nil, errors.New("error: find many items failed")
		}
		results = append(results, &item.ItemShowCase{
			ID:       result.ID.Hex(),
			Title:    result.Title,
			Price:    result.Price,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) CountItems(c context.Context, filter primitive.D) (int64, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: CountItems: %s\n", err.Error())
		return -1, errors.New("error: count items failed")
	}

	return count, nil
}

func (r *itemRepository) UpdateItem(c context.Context, itemID string, req primitive.M) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	if _, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemID)}, bson.M{"$set": req}); err != nil {
		log.Printf("Error: UpdateItem: %s\n", err.Error())
		return errors.New("error: update item failed")
	}

	return nil
}

func (r *itemRepository) UpdateItemStatus(c context.Context, itemID string, isActive bool) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	if _, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(itemID)}, bson.M{"$set": bson.M{"usage_status": isActive}},
	); err != nil {
		log.Printf("Error: UpdateItemStatus: %s\n", err.Error())
		return errors.New("error: update item status failed")
	}

	return nil
}
