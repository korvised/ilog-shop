package itemHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
)

type (
	ItemHttpHandlerService interface {
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecases.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecases.ItemUsecaseService) ItemHttpHandlerService {
	return &itemHttpHandler{cfg, itemUsecase}
}
