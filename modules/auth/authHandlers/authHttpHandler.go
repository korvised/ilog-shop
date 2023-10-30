package authHandlers

import (
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth/authUsecases"
)

type (
	AuthHttpHandlerService interface {
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authUsecases.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authUsecases.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg, authUsecase}
}
