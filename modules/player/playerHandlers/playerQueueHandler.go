package playerHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
)

type (
	PlayerQueueHandlerService interface {
	}

	playerQueueHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerQueueHandler(cfg *config.Config, playerUsecase playerUsecases.PlayerUsecaseService) PlayerQueueHandlerService {
	return &playerQueueHandler{cfg, playerUsecase}
}
