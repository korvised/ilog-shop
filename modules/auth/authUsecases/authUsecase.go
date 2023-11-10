package authUsecases

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	"github.com/korvised/ilog-shop/modules/auth/authRepositories"
)

type (
	AuthUsecaseService interface {
		Login(c context.Context, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error)
		Logout(c context.Context, req *auth.PlayerLogoutReq) error
		RefreshToken(c context.Context, req *auth.RefreshTokenReq) (*auth.ProfileInterceptor, error)
		GetCredentialByAccessToken(c context.Context, req *authPb.CredentialReq) (*authPb.CredentialRes, error)
		GetRoleCount(c context.Context) (*authPb.RolesCountRes, error)
	}

	authUsecase struct {
		cfg            *config.Config
		authRepository authRepositories.AuthRepositoryService
	}
)

func NewAuthUsecase(cfg *config.Config, authRepository authRepositories.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{cfg, authRepository}
}
