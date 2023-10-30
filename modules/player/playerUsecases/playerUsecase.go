package playerUsecases

import "github.com/korvised/ilog-shop/modules/player/playerRepositories"

type (
	PlayerUsecaseService interface {
	}

	playerUsecase struct {
		playerRepository playerRepositories.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepositories.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository}
}
