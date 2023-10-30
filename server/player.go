package server

import (
	"github.com/korvised/ilog-shop/modules/player/playerHandlers"
	"github.com/korvised/ilog-shop/modules/player/playerRepositories"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
)

func (s *server) playerService() {
	repo := playerRepositories.NewPlayerRepository(s.db)
	usecase := playerUsecases.NewPlayerUsecase(repo)
	httpHandler := playerHandlers.NewPlayerHttpHandler(s.cfg, usecase)
	grpcHandler := playerHandlers.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandlers.NewPlayerQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	player := s.app.Group("/player_v1")

	// Health check
	_ = player
}
