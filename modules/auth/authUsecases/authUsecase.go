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
	"log"
)

type (
	AuthUsecaseService interface {
		Login(c context.Context, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error)
		Logout(c context.Context, req *auth.PlayerLogoutReq) error
		RefreshToken(c context.Context, req *auth.RefreshTokenReq) (*auth.ProfileInterceptor, error)
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

	credentialId, err := u.authRepository.InsertCredential(c, &auth.Credential{
		PlayerID:     profile.GetId(),
		RoleCode:     roleCode,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})

	credential, err := u.authRepository.FindOneCredential(c, credentialId.Hex())
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

func (u *authUsecase) Logout(c context.Context, req *auth.PlayerLogoutReq) error {
	credential, err := u.authRepository.FindOneCredentialByAccessToken(c, req.AccessToken)
	if err != nil {
		return err
	}

	if err = u.authRepository.DeleteCredential(c, credential.ID.Hex()); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) RefreshToken(c context.Context, req *auth.RefreshTokenReq) (*auth.ProfileInterceptor, error) {
	claims, err := jwtauth.ParseToken(u.cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: RefreshToken: %v \n", err)
		return nil, err
	}

	profile, err := u.authRepository.FindOnePlayerProfileToRefresh(c, u.cfg.Grpc.PlayerUrl, &playerPb.FindOnePlayerProfileToRefreshReq{
		PlayerId: claims.PlayerID,
	})
	if err != nil {
		return nil, err
	}

	roleCode := int(profile.GetRoleCode())

	accessToken := jwtauth.NewAccessToken(u.cfg.Jwt.AccessSecretKey, u.cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerID: profile.GetId(),
		RoleCode: roleCode,
	}).SingToken()

	refreshToken := jwtauth.ReloadToken(u.cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtauth.Claims{
		PlayerID: profile.GetId(),
		RoleCode: roleCode,
	})

	credential, err := u.authRepository.FindOneCredentialByRefreshToken(c, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if err = u.authRepository.UpdateCredential(c, &auth.UpdateCredentialReq{
		ID:           credential.ID.Hex(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err = u.authRepository.FindOneCredential(c, credential.ID.Hex())
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
