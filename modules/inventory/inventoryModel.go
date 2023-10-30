package inventory

import (
	"github.com/korvised/ilog-shop/modules/item"
	"github.com/korvised/ilog-shop/modules/models"
)

type (
	UpdateInventoryReq struct {
		PlayerID string `json:"player_id" validate:"required,max=64"`
		ItemID   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInventory struct {
		InventoryID string `json:"inventory_id"`
		*item.ItemShowCase
	}

	PlayerInventory struct {
		PlayerID string `json:"player_id"`
		*models.PaginateRes
	}
)
