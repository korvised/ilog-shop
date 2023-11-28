package inventoryHandlers

import (
	"context"
	inventoryPb "github.com/korvised/ilog-shop/modules/inventory/inventoryPb"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryUsecases"
)

type (
	inventoryGrpcHandler struct {
		inventoryPb.UnimplementedInventoryGrpcServiceServer
		inventoryUsecase inventoryUsecases.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecases.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryGrpcHandler) IsAvailableToSell(
	ctx context.Context,
	req *inventoryPb.IsAvailableToSellReq,
) (*inventoryPb.IsAvailableToSellRes, error) {
	return nil, nil
}
