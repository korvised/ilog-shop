package itemUsecases

import "github.com/korvised/ilog-shop/modules/item/itemRepositories"

type (
	ItemUsecaseService interface {
	}

	itemUsecase struct {
		itemRepository itemRepositories.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepository itemRepositories.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{itemRepository}
}
