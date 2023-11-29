package paymentUsecases

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/modules/inventory"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/queue"
	"github.com/korvised/ilog-shop/pkg/utils"
	"log"
)

func (u *paymentUsecase) BuyOrSellConsumer(c context.Context, key string, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(c)
	if err != nil {
		resCh <- nil
		return
	}
	defer consumer.Close()

	log.Println("Start BuyOrSellConsumer ...")

	select {
	case err = <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer failed: ", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			_ = u.UpsertOffset(c, msg.Offset+1)

			req := new(payment.PaymentTransferRes)

			if err = queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("BuyOrSellConsumer | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
		}
	}
}

func (u *paymentUsecase) BuyItem(c context.Context, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.GetPlayerItems(c, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, obj := range req.Items {
		_ = u.paymentRepository.DockedPlayerMoney(c, &player.CreatePlayerTransactionReq{
			PlayerID: playerID,
			Amount:   -obj.Price,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(c, "buy", resCh)

		res := <-resCh
		if res != nil {
			utils.Debug(res, "stage1 res:")
			stage1 = append(stage1, &payment.PaymentTransferRes{
				InventoryID:   "",
				TransactionID: res.TransactionID,
				PlayerID:      playerID,
				ItemID:        obj.ItemID,
				Amount:        obj.Price,
				Error:         res.Error,
			})
		}
	}

	// Rollback transaction
	for _, obj := range stage1 {
		if obj.Error != "" {
			log.Println("BuyOrSellConsumer stage1 error: ", obj)
			for _, s1 := range stage1 {
				_ = u.paymentRepository.RollbackTransaction(c, &player.RollbackPlayerTransactionReq{
					TransactionID: s1.TransactionID,
				})
			}

			return nil, errors.New("error: buy item failed")
		}
	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		_ = u.paymentRepository.AddPlayItem(c, &inventory.UpdateInventoryReq{
			PlayerID: playerID,
			ItemID:   s1.ItemID,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(c, "buy", resCh)

		res := <-resCh
		if res != nil {
			utils.Debug(res, "stage2 res: ")
			stage2 = append(stage2, &payment.PaymentTransferRes{
				InventoryID:   res.InventoryID,
				TransactionID: s1.TransactionID,
				PlayerID:      playerID,
				ItemID:        s1.ItemID,
				Amount:        s1.Amount,
				Error:         res.Error,
			})
		}
	}

	// Rollback inventory
	for _, obj := range stage2 {
		if obj.Error != "" {
			log.Println("BuyOrSellConsumer stage2 error: ", obj)
			for _, s1 := range stage1 {
				_ = u.paymentRepository.RollbackTransaction(c, &player.RollbackPlayerTransactionReq{
					TransactionID: s1.TransactionID,
				})
			}

			for _, s2 := range stage2 {
				_ = u.paymentRepository.RollbackAddPlayItem(c, &inventory.RollbackPlayerInventoryReq{
					InventoryID: s2.InventoryID,
				})
			}

			return nil, errors.New("error: buy item failed")
		}
	}

	return stage2, nil
}

func (u *paymentUsecase) SellItem(c context.Context, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.GetPlayerItems(c, req.Items); err != nil {
		return nil, err
	}

	return nil, nil
}
