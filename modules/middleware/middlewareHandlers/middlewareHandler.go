package middlewareHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareUsecases"
	"github.com/korvised/ilog-shop/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type (
	MiddlewareHandlerService interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
		Roles(next echo.HandlerFunc, expectedRoles []int) echo.HandlerFunc
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecases.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecases.MiddlewareUsecaseService) MiddlewareHandlerService {
	return &middlewareHandler{cfg, middlewareUsecase}
}

func (h *middlewareHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization == "" {
			return response.ErrResponse(c, http.StatusUnauthorized, "error: authorization")
		}

		authorizationParts := strings.Split(authorization, " ")
		if len(authorizationParts) != 2 {
			return response.ErrResponse(c, http.StatusUnauthorized, "error: authorization")
		}

		accessToken := authorizationParts[1]

		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, accessToken)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) Roles(next echo.HandlerFunc, expectedRoles []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.Roles(c, expectedRoles)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}
