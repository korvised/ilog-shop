package itemUsecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/item"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/item/itemRepositories"
	"github.com/korvised/ilog-shop/modules/models"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemUsecaseService interface {
		CreateItem(c context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		GetItem(c context.Context, itemID string) (*item.ItemShowCase, error)
		GetItems(c context.Context, req *item.ItemSearchReq) (*models.PaginateRes, error)
		GetItemsInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		UpdateItem(c context.Context, itemID string, req *item.ItemUpdateReq) (*item.ItemShowCase, error)
		UpdateItemStatus(c context.Context, itemID string, isActive bool) (*item.ItemShowCase, error)
	}

	itemUsecase struct {
		cfg            *config.Config
		itemRepository itemRepositories.ItemRepositoryService
	}
)

func NewItemUsecase(cfg *config.Config, itemRepository itemRepositories.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{cfg, itemRepository}
}

func (u *itemUsecase) CreateItem(c context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error) {
	if !u.itemRepository.IsUniqueItem(c, req.Title) {
		return nil, errors.New("error: item already exists")
	}

	payload := &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		UsageStatus: true,
		ImageUrl:    req.ImageUrl,
		CreatedAt:   utils.LocalTime(),
		UpdatedAt:   utils.LocalTime(),
	}

	itemID, err := u.itemRepository.InsertOneItem(c, payload)
	if err != nil {
		return nil, err
	}

	return u.GetItem(c, itemID.Hex())
}

func (u *itemUsecase) GetItem(c context.Context, itemID string) (*item.ItemShowCase, error) {
	result, err := u.itemRepository.FindOneItem(c, itemID)
	if err != nil {
		return nil, err
	}

	return &item.ItemShowCase{
		ID:       result.ID.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}

func (u *itemUsecase) GetItems(c context.Context, req *item.ItemSearchReq) (*models.PaginateRes, error) {
	filter := bson.D{}
	countFilter := bson.D{}
	opts := make([]*options.FindOptions, 0)

	// Filter
	if req.Start != "" {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}

	if req.Title != "" {
		filter = append(filter, bson.E{Key: "title", Value: primitive.Regex{Pattern: req.Title, Options: "i"}})
		countFilter = append(countFilter, bson.E{Key: "title", Value: primitive.Regex{Pattern: req.Title, Options: "i"}})
	}

	filter = append(filter, bson.E{Key: "usage_status", Value: true})
	countFilter = append(filter, bson.E{Key: "usage_status", Value: true})

	// Options
	opts = append(opts, options.Find().SetSort(bson.D{{"_id", -1}}))
	opts = append(opts, options.Find().SetLimit(int64(req.Limit)))

	// Count
	total, err := u.itemRepository.CountItems(c, countFilter)
	if err != nil {
		return nil, err
	}

	// Find
	results, err := u.itemRepository.FindManyItems(c, filter, opts)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*item.ItemShowCase, 0),
			Limit: req.Limit,
			Total: total,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s?title=%s&limit=%d", u.cfg.Paginate.ItemNextPageBaseURL, req.Title, req.Limit),
			},
			Next: models.NextPaginate{
				Href:  "",
				Start: "",
			},
		}, nil
	}

	// Paginate
	return &models.PaginateRes{
		Data:  results,
		Limit: req.Limit,
		Total: total,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s?title=%s&limit=%d", u.cfg.Paginate.ItemNextPageBaseURL, req.Title, req.Limit),
		},
		Next: models.NextPaginate{
			Href: fmt.Sprintf(
				"%s?title=%s&start=%s&limit=%d",
				u.cfg.Paginate.ItemNextPageBaseURL,
				req.Title, results[len(results)-1].ID,
				req.Limit,
			),
			Start: results[len(results)-1].ID,
		},
	}, nil
}

func (u *itemUsecase) GetItemsInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	filter := bson.D{}

	objIds := make([]primitive.ObjectID, 0)
	for _, id := range req.Ids {
		objIds = append(objIds, utils.ConvertToObjectId(id))
	}

	filter = append(filter, bson.E{Key: "_id", Value: bson.D{{"$in", objIds}}})
	filter = append(filter, bson.E{Key: "usage_status", Value: true})

	results, err := u.itemRepository.FindManyItems(c, filter, nil)
	if err != nil {
		return nil, err
	}

	items := make([]*itemPb.Item, 0)
	for _, obj := range results {
		items = append(items, &itemPb.Item{
			Id:       obj.ID,
			Title:    obj.Title,
			Price:    obj.Price,
			Damage:   int32(obj.Damage),
			ImageUrl: obj.ImageUrl,
		})
	}

	result := &itemPb.FindItemsInIdsRes{
		Items: items,
	}

	return result, nil
}

func (u *itemUsecase) UpdateItem(c context.Context, itemID string, req *item.ItemUpdateReq) (*item.ItemShowCase, error) {
	payload := bson.M{}

	if req.Title != "" {
		payload["title"] = req.Title
	}

	if req.Price >= 0 {
		payload["price"] = req.Price
	}

	if req.Damage != 0 {
		payload["damage"] = req.Damage
	}

	if req.ImageUrl != "" {
		payload["image_url"] = req.ImageUrl
	}

	payload["updated_at"] = utils.LocalTime()

	if err := u.itemRepository.UpdateItem(c, itemID, payload); err != nil {
		return nil, err
	}

	return u.GetItem(c, itemID)
}

func (u *itemUsecase) UpdateItemStatus(c context.Context, itemID string, isActive bool) (*item.ItemShowCase, error) {
	if err := u.itemRepository.UpdateItemStatus(c, itemID, isActive); err != nil {
		return nil, err
	}

	return u.GetItem(c, itemID)
}
