package playerHandlers

import (
	"context"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
)

type (
	playerGrpcHandler struct {
		playerPb.UnimplementedPlayerGrpcServiceServer
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecases.PlayerUsecaseService) *playerGrpcHandler {
	return &playerGrpcHandler{
		playerUsecase: playerUsecase,
	}
}

func (g *playerGrpcHandler) CredentialSearch(
	ctx context.Context,
	req *playerPb.CredentialSearchReq,
) (*playerPb.PlayerProfile, error) {
	return g.playerUsecase.GetPlayerCredential(ctx, req.Email, req.Password)
}

func (g *playerGrpcHandler) FindOnePlayerProfileToRefresh(
	ctx context.Context,
	req *playerPb.FindOnePlayerProfileToRefreshReq,
) (*playerPb.PlayerProfile, error) {
	return nil, nil
}

func (g *playerGrpcHandler) GetPlayerSavingAccount(
	ctx context.Context,
	req *playerPb.GetPlayerSavingAccountReq,
) (*playerPb.GetPlayerSavingAccountRes, error) {
	return nil, nil
}
