package inventoryUsecases

import (
	"context"
	"github.com/korvised/ilog-shop/modules/inventory"
	"github.com/korvised/ilog-shop/modules/payment"
)

func (u *inventoryUsecase) AddPlayerItemRes(c context.Context, req *inventory.UpdateInventoryReq) {
	inventoryID, err := u.inventoryRepository.InsertOnePlayerItem(c, &inventory.Inventory{
		PlayerID: req.PlayerID,
		ItemID:   req.ItemID,
	})
	if err != nil {
		_ = u.inventoryRepository.AddPlayerItemRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         err.Error(),
		})

		return
	}

	_ = u.inventoryRepository.AddPlayerItemRes(c, &payment.PaymentTransferRes{
		InventoryID:   inventoryID.Hex(),
		TransactionID: "",
		PlayerID:      req.PlayerID,
		ItemID:        req.ItemID,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RemovePlayerItemRes(c context.Context, req *inventory.UpdateInventoryReq) {
	if err := u.inventoryRepository.FindOnePlayerItem(c, req.PlayerID, req.ItemID); err != nil {
		_ = u.inventoryRepository.RemovePlayerItemRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         err.Error(),
		})

		return
	}

	if err := u.inventoryRepository.DeleteOnePlayerItem(c, req.PlayerID, req.ItemID); err != nil {
		_ = u.inventoryRepository.RemovePlayerItemRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         err.Error(),
		})

		return
	}

	_ = u.inventoryRepository.RemovePlayerItemRes(c, &payment.PaymentTransferRes{
		InventoryID:   "",
		TransactionID: "",
		PlayerID:      req.PlayerID,
		ItemID:        req.ItemID,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RollbackAddPlayerItem(c context.Context, req *inventory.RollbackPlayerInventoryReq) {
	_ = u.inventoryRepository.DeleteOneInventory(c, req.InventoryID)
}

func (u *inventoryUsecase) RollbackRemovePlayerItem(c context.Context, req *inventory.RollbackPlayerInventoryReq) {
	_, _ = u.inventoryRepository.InsertOnePlayerItem(c, &inventory.Inventory{
		PlayerID: req.PlayerID,
		ItemID:   req.ItemID,
	})
}
