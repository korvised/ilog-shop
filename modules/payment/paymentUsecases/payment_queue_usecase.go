package paymentUsecases

import (
	"context"
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
			utils.Debug(res)
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

	for _, obj := range stage1 {
		if obj.Error != "" {
			log.Println("BuyOrSellConsumer error: ", obj)
			for _, s1 := range stage1 {
				_ = u.paymentRepository.RollbackTransaction(c, &player.RollbackPlayerTransactionReq{
					TransactionID: s1.TransactionID,
				})
			}
		}
	}

	return stage1, nil
}

func (u *paymentUsecase) SellItem(c context.Context, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.GetPlayerItems(c, req.Items); err != nil {
		return nil, err
	}

	return nil, nil
}
