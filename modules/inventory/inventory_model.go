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

	ItemInInventory struct {
		InventoryID string `json:"inventory_id"`
		PlayerID    string `json:"player_id"`
		*item.ItemShowCase
	}

	InventorySearchReq struct {
		models.PaginateReq
	}
)
