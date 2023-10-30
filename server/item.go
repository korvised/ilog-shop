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

	router := s.app.Group("/api/v1/item")

	// Health check
	router.GET("", s.healthCheckService)
}
