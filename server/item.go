package server

import (
	"github.com/korvised/ilog-shop/modules/item/itemHandlers"
	"github.com/korvised/ilog-shop/modules/item/itemRepositories"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
)

func (s *server) itemService() {
	repo := itemRepositories.NewItemRepository(s.db)
	usecase := itemUsecases.NewItemUsecase(repo)
	httpHandler := itemHandlers.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandlers.NewItemGrpcHandler(usecase)

	_ = httpHandler
	_ = grpcHandler

	item := s.app.Group("/item_v1")

	// Health check
	_ = item
}
