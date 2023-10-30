package server

import (
	"github.com/korvised/ilog-shop/modules/inventory/inventoryHandlers"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryRepositories"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryUsecases"
)

func (s *server) inventoryService() {
	repo := inventoryRepositories.NewInventoryRepository(s.db)
	usecase := inventoryUsecases.NewInventoryUsecase(repo)
	httpHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandlers.NewInventoryGrpcHandler(usecase)
	queueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	inventory := s.app.Group("/inventory_v1")

	// Health check
	_ = inventory
}
