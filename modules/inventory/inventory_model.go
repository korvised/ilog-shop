package inventory

import (
	"github.com/korvised/ilog-shop/modules/item"
	"github.com/korvised/ilog-shop/modules/models"
)

type (
	UpdateInventoryReq struct {
		PlayerID string `json:"player_id"`
		ItemID   string `json:"item_id"`
	}

	ItemInInventory struct {
		InventoryID string `json:"inventory_id"`
		PlayerID    string `json:"player_id"`
		*item.ItemShowCase
	}

	InventorySearchReq struct {
		models.PaginateReq
	}

	RollbackPlayerInventoryReq struct {
		InventoryID string `json:"inventory_id"`
		PlayerID    string `json:"player_id"`
		ItemID      string `json:"item_id"`
	}
)
