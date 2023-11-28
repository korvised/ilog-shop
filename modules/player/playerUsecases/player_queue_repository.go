package playerUsecases

import (
	"context"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/utils"
	"log"
	"math"
)

func (u *playerUsecase) DockedPlayerMoneyRes(c context.Context, req *player.CreatePlayerTransactionReq) {
	// Get saving account
	savingAccount, err := u.playerRepository.FindOnePlayerSavingAccount(c, req.PlayerID)
	if err != nil {
		log.Printf("Error: DockedPlayerMoneyRes: %s\n", err)
		u.playerRepository.DockedPlayerMoneyRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})
		return
	}

	if savingAccount.Balance < math.Abs(req.Amount) {
		log.Printf("Error: DockedPlayerMoneyRes failed: %s", "not enough money")
		u.playerRepository.DockedPlayerMoneyRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        "",
			Amount:        req.Amount,
			Error:         "error: not enough money",
		})
		return
	}

	// Insert one player transaction
	transactionID, err := u.playerRepository.InsertOnePlayerTransaction(c, &player.PlayerTransaction{
		PlayerID:  req.PlayerID,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	})
	if err != nil {
		u.playerRepository.DockedPlayerMoneyRes(c, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})
		return
	}

	u.playerRepository.DockedPlayerMoneyRes(c, &payment.PaymentTransferRes{
		InventoryID:   "",
		TransactionID: transactionID.Hex(),
		PlayerID:      req.PlayerID,
		ItemID:        "",
		Amount:        req.Amount,
		Error:         "",
	})
}
