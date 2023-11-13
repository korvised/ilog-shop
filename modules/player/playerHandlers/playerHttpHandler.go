package playerHandlers

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
	"github.com/korvised/ilog-shop/pkg/request"
	"github.com/korvised/ilog-shop/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	PlayerHttpHandlerService interface {
		AddPlayerMoney(c echo.Context) error
		GetPlayerProfile(c echo.Context) error
		GetPlayerSavingAccount(c echo.Context) error
		Register(c echo.Context) error
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerHttpHandler(cfg *config.Config, playerUsecase playerUsecases.PlayerUsecaseService) PlayerHttpHandlerService {
	return &playerHttpHandler{cfg, playerUsecase}
}

func (h *playerHttpHandler) Register(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)
	req := new(player.CreatePlayerReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.CreatePlayer(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *playerHttpHandler) GetPlayerProfile(c echo.Context) error {
	ctx := context.Background()

	playerId := c.Param("player_id")
	if playerId == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "error: player id is required")
	}

	res, err := h.playerUsecase.GetPlayerProfile(ctx, playerId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) AddPlayerMoney(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(player.CreatePlayerTransactionReq)
	req.PlayerID = c.Get(middleware.PlayerID).(string)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.AddPlayerMoney(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) GetPlayerSavingAccount(c echo.Context) error {
	ctx := context.Background()

	playerId := c.Get(middleware.PlayerID).(string)

	res, err := h.playerUsecase.GetPlayerSavingAccount(ctx, playerId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
