package inventoryUsecases

import (
	"context"
	"fmt"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/inventory"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryRepositories"
	"github.com/korvised/ilog-shop/modules/item"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/models"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type (
	InventoryUsecaseService interface {
		GetPlayerItems(c context.Context, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
	}

	inventoryUsecase struct {
		cfg                 *config.Config
		inventoryRepository inventoryRepositories.InventoryRepositoryService
	}
)

func NewInventoryUsecase(cfg *config.Config, inventoryRepository inventoryRepositories.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{cfg, inventoryRepository}
}

func (u *inventoryUsecase) GetPlayerItems(c context.Context, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
	// Options
	opts := make([]*options.FindOptions, 0)
	opts = append(opts, options.Find().SetSort(bson.D{{"_id", -1}}))
	opts = append(opts, options.Find().SetLimit(int64(req.Limit)))

	// Filter
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "player_id", Value: playerID})

	// Count
	total, err := u.inventoryRepository.CountPlayerItems(c, filter)

	if req.Start != "" {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}

	// Find
	inventoryData, err := u.inventoryRepository.FindPlayItems(c, filter, opts)
	if err != nil {
		return nil, err
	}

	nextPaginate := models.NextPaginate{
		Href:  "",
		Start: "",
	}

	results := make([]*inventory.ItemInInventory, 0)
	if len(inventoryData) > 0 {
		nextPaginate.Href = fmt.Sprintf(
			"%s/%s?start=%s&limit=%d",
			u.cfg.Paginate.InventoryNextPageBaseURL,
			playerID,
			inventoryData[len(inventoryData)-1].ID.Hex(),
			req.Limit,
		)
		nextPaginate.Start = inventoryData[len(inventoryData)-1].ID.Hex()

		itemData, err := u.inventoryRepository.FindItemInIds(c, &itemPb.FindItemsInIdsReq{
			Ids: func() []string {
				itemIds := make([]string, 0)
				for _, obj := range inventoryData {
					itemIds = append(itemIds, obj.ItemID)
				}

				return itemIds
			}(),
		})
		if err != nil {
			log.Printf("Error: GetPlayerItems: %v \n", err)
			return nil, fmt.Errorf("error: item not found")
		}

		itemMaps := make(map[string]*item.ItemShowCase)
		for _, obj := range itemData.Items {
			itemMaps[obj.Id] = &item.ItemShowCase{
				ID:       obj.Id,
				Title:    obj.Title,
				Price:    obj.Price,
				Damage:   int(obj.Damage),
				ImageUrl: obj.ImageUrl,
			}
		}

		for _, obj := range inventoryData {
			results = append(results, &inventory.ItemInInventory{
				InventoryID:  obj.ID.Hex(),
				PlayerID:     obj.PlayerID,
				ItemShowCase: itemMaps[obj.ItemID],
			})
		}
	}

	return &models.PaginateRes{
		Data:  results,
		Limit: req.Limit,
		Total: total,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", u.cfg.Paginate.InventoryNextPageBaseURL, playerID, req.Limit),
		},
		Next: nextPaginate,
	}, nil
}
