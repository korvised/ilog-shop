package authHandlers

import "github.com/korvised/ilog-shop/modules/auth/authUsecases"

type (
	authGrpcHandler struct {
		authUsecase authUsecases.AuthUsecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecases.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{authUsecase}
}
