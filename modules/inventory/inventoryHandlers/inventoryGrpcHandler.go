package inventoryHandlers

import "github.com/korvised/ilog-shop/modules/inventory/inventoryUsecases"

type (
	inventoryGrpcHandler struct {
		inventoryUsecase inventoryUsecases.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecases.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUsecase}
}
