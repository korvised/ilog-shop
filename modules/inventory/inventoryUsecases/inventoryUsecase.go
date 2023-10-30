package inventoryUsecases

import "github.com/korvised/ilog-shop/modules/inventory/inventoryRepositories"

type (
	InventoryUsecaseService interface {
	}

	inventoryUsecase struct {
		inventoryRepository inventoryRepositories.InventoryRepositoryService
	}
)

func NewInventoryUsecase(inventoryRepository inventoryRepositories.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{inventoryRepository}
}
