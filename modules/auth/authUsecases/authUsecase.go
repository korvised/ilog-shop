package authUsecases

import "github.com/korvised/ilog-shop/modules/auth/authRepositories"

type (
	AuthUsecaseService interface {
	}

	authUsecase struct {
		authRepository authRepositories.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepository authRepositories.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{authRepository}
}
