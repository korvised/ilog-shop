package middlewareHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareUsecases"
)

type (
	MiddlewareHandlerService interface {
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecases.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecases.MiddlewareUsecaseService) MiddlewareHandlerService {
	return &middlewareHandler{cfg, middlewareUsecase}
}
