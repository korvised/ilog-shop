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
