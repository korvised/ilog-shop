package itemHandlers

import "github.com/korvised/ilog-shop/modules/item/itemUsecases"

type (
	itemGrpcHandler struct {
		itemUsecase itemUsecases.ItemUsecaseService
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecases.ItemUsecaseService) *itemGrpcHandler {
	return &itemGrpcHandler{itemUsecase}
}
