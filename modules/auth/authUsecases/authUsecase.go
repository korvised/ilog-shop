package authUsecases

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth"
	"github.com/korvised/ilog-shop/modules/auth/authRepositories"
	"github.com/korvised/ilog-shop/modules/player"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"github.com/korvised/ilog-shop/pkg/utils"
)

type (
	AuthUsecaseService interface {
		Login(c context.Context, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error)
	}

	authUsecase struct {
		cfg            *config.Config
		authRepository authRepositories.AuthRepositoryService
	}
)

func NewAuthUsecase(cfg *config.Config, authRepository authRepositories.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{cfg, authRepository}
}

func (u *authUsecase) Login(c context.Context, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error) {

	profile, err := u.authRepository.CredentialSearch(c, u.cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err

	}

	roleCode := int(profile.GetRoleCode())

	accessToken := jwtauth.NewAccessToken(u.cfg.Jwt.AccessSecretKey, u.cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerID: profile.GetId(),
		RoleCode: roleCode,
	}).SingToken()

	refreshToken := jwtauth.NewRefreshToken(u.cfg.Jwt.RefreshSecretKey, u.cfg.Jwt.RefreshDuration, &jwtauth.Claims{
		PlayerID: profile.GetId(),
		RoleCode: roleCode,
	}).SingToken()

	credentialId, err := u.authRepository.InsertOnePlayerCredential(c, &auth.Credential{
		PlayerID:     profile.GetId(),
		RoleCode:     roleCode,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})

	credential, err := u.authRepository.FindOnePlayerCredential(c, credentialId.Hex())
	if err != nil {
		return nil, err
	}

	return &auth.ProfileInterceptor{
		PlayerProfile: &player.PlayerProfile{
			ID:        profile.GetId(),
			Email:     profile.GetEmail(),
			Username:  profile.GetUsername(),
			CreatedAt: utils.ConvertStringTimeToTime(profile.GetCreatedAt()),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.GetUpdatedAt()),
		},
		Credential: &auth.CredentialRes{
			ID:           credential.ID.Hex(),
			PlayerID:     credential.PlayerID,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt,
			UpdatedAt:    credential.UpdatedAt,
		},
	}, nil
}
