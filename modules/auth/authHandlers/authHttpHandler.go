package authHandlers

import (
	"context"
	"fmt"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth"
	"github.com/korvised/ilog-shop/modules/auth/authUsecases"
	"github.com/korvised/ilog-shop/pkg/request"
	"github.com/korvised/ilog-shop/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	AuthHttpHandlerService interface {
		Login(c echo.Context) error
		Logout(c echo.Context) error
		RefreshToken(c echo.Context) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authUsecases.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authUsecases.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg, authUsecase}
}

func (h *authHttpHandler) Login(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	req := new(auth.PlayerLoginReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Login(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *authHttpHandler) Logout(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	req := new(auth.PlayerLogoutReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.authUsecase.Logout(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, &response.MsgResponse{
		Message: "Logout success",
	})
}

func (h *authHttpHandler) RefreshToken(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	req := new(auth.RefreshTokenReq)

	fmt.Println(req)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.RefreshToken(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
