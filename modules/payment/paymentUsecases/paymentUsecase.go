package paymentUsecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/item"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/modules/payment/paymentRepositories"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/queue"
	"log"
)

type (
	PaymentUsecaseService interface {
		GetPlayerItems(c context.Context, req []*payment.ItemServiceReqDatum) error
		GetOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
		PaymentConsumer(c context.Context) (sarama.PartitionConsumer, error)
		BuyOrSellConsumer(c context.Context, key string, resCh chan<- *payment.PaymentTransferRes)
		BuyItem(c context.Context, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		SellItem(c context.Context, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
	}

	paymentUsecase struct {
		cfg               *config.Config
		paymentRepository paymentRepositories.PaymentRepositoryService
	}
)

func NewPaymentUsecase(cfg *config.Config, paymentRepository paymentRepositories.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{cfg, paymentRepository}
}

func (u *paymentUsecase) GetPlayerItems(c context.Context, req []*payment.ItemServiceReqDatum) error {
	setIds := make(map[string]bool)
	for _, obj := range req {
		if !setIds[obj.ItemID] {
			setIds[obj.ItemID] = true
		}
	}

	itemData, err := u.paymentRepository.FindItemInIds(c, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for key := range setIds {
				itemIds = append(itemIds, key)
			}

			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: GetPlayerItems: %v \n", err)
		return fmt.Errorf("error: item not found")
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

	for i := range req {
		if _, ok := itemMaps[req[i].ItemID]; !ok {
			log.Printf("Error: GetPlayerItems: item not found: %v \n", req[i].ItemID)
			return errors.New("error: item not found")
		}

		req[i].Price = itemMaps[req[i].ItemID].Price
	}

	return nil
}

func (u *paymentUsecase) GetOffset(c context.Context) (int64, error) {
	return u.paymentRepository.FindOffset(c)
}

func (u *paymentUsecase) UpsertOffset(c context.Context, offset int64) error {
	return u.paymentRepository.UpsertOffset(c, offset)
}

func (u *paymentUsecase) PaymentConsumer(c context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{u.cfg.Kafka.Url}, u.cfg.Kafka.ApiKey, u.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepository.FindOffset(c)
	if err != nil {
		return nil, err
	}
	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		log.Printf("Error: PaymentConsumer 1: %v \n", err)
		consumer, err = worker.ConsumePartition("payment", 0, offset)
		if err != nil {
			log.Printf("Error: PaymentConsumer 2 : %v \n", err)
			return nil, err
		}
	}

	return consumer, nil
}

func (u *paymentUsecase) BuyOrSellConsumer(c context.Context, key string, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(c)
	if err != nil {
		resCh <- nil
		return
	}
	defer consumer.Close()

	log.Println("Start BuyOrSellConsumer ...")

	select {
	case err := <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer failed: ", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpsertOffset(c, msg.Offset+1)

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
			log.Println("BuyOrSellConsumer res: ", res)
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
