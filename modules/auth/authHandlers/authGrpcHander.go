package authHandlers

import (
	"context"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	"github.com/korvised/ilog-shop/modules/auth/authUsecases"
)

type (
	authGrpcHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase authUsecases.AuthUsecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecases.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{
		authUsecase: authUsecase,
	}
}

func (g *authGrpcHandler) GetCredential(ctx context.Context, req *authPb.CredentialReq) (*authPb.CredentialRes, error) {
	return g.authUsecase.GetCredentialByAccessToken(ctx, req)
}

func (g *authGrpcHandler) GetRolesCount(ctx context.Context, _ *authPb.RolesCountReq) (*authPb.RolesCountRes, error) {
	return g.authUsecase.GetRoleCount(ctx)
}
