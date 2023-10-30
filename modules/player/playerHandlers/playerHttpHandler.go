package playerHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
)

type (
	PlayerHttpHandlerService interface {
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerHttpHandler(cfg *config.Config, playerUsecase playerUsecases.PlayerUsecaseService) PlayerHttpHandlerService {
	return &playerHttpHandler{cfg, playerUsecase}
}
