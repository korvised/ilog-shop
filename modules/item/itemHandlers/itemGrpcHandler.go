package itemHandlers

import (
	"context"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
)

type (
	itemGrpcHandler struct {
		itemPb.UnimplementedItemGrpcServiceServer
		itemUsecase itemUsecases.ItemUsecaseService
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecases.ItemUsecaseService) *itemGrpcHandler {
	return &itemGrpcHandler{
		itemUsecase: itemUsecase,
	}
}

func (g *itemGrpcHandler) FindItemInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return nil, nil
}
