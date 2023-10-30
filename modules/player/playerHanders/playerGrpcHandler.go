package playerHanders

import "github.com/korvised/ilog-shop/modules/player/playerUsecases"

type (
	playerGrpcHandlerService struct {
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecases.PlayerUsecaseService) *playerGrpcHandlerService {
	return &playerGrpcHandlerService{playerUsecase}
}
