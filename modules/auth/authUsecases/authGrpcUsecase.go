package authUsecases

import (
	"context"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
)

func (u *authUsecase) GetCredentialByAccessToken(c context.Context, req *authPb.CredentialReq) (*authPb.CredentialRes, error) {
	_, err := u.authRepository.FindOneCredentialByAccessToken(c, req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return &authPb.CredentialRes{
		IsValid: true,
	}, nil
}

func (u *authUsecase) GetRoleCount(c context.Context) (*authPb.RolesCountRes, error) {
	count, err := u.authRepository.FindRoleCount(c)
	if err != nil {
		return nil, err
	}

	return &authPb.RolesCountRes{
		Count: count,
	}, nil
}
