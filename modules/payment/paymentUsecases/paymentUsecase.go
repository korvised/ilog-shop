package paymentUsecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/korvised/ilog-shop/modules/item"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/modules/payment/paymentRepositories"
	"log"
)

type (
	PaymentUsecaseService interface {
		GetPlayerItems(c context.Context, req []*payment.ItemServiceReqDatum) error
		GetOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
	}

	paymentUsecase struct {
		paymentRepository paymentRepositories.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepository paymentRepositories.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{paymentRepository}
}

func (u *paymentUsecase) GetPlayerItems(c context.Context, req []*payment.ItemServiceReqDatum) error {
	setIds := make(map[string]bool)
	for _, obj := range req {
		if !setIds[obj.ItemId] {
			setIds[obj.ItemId] = true
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
		if _, ok := itemMaps[req[i].ItemId]; !ok {
			log.Printf("Error: GetPlayerItems: item not found: %v \n", req[i].ItemId)
			return errors.New("error: item not found")
		}

		req[i].Price = itemMaps[req[i].ItemId].Price
	}

	return nil
}

func (u *paymentUsecase) GetOffset(c context.Context) (int64, error) {
	return u.paymentRepository.FindOffset(c)
}

func (u *paymentUsecase) UpsertOffset(c context.Context, offset int64) error {
	return u.paymentRepository.UpsertOffset(c, offset)
}
