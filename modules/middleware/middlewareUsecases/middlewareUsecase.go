package middlewareUsecases

import (
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareRepositories"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"github.com/korvised/ilog-shop/pkg/rbac"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, accessToken string) (echo.Context, error)
		Roles(c echo.Context, expectedRoles []int) (echo.Context, error)
	}

	middlewareUsecase struct {
		cfg                  *config.Config
		middlewareRepository middlewareRepositories.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(cfg *config.Config, middlewareRepository middlewareRepositories.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{cfg, middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, accessToken string) (echo.Context, error) {
	ctx := c.Request().Context()

	claims, err := jwtauth.ParseToken(u.cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err = u.middlewareRepository.FindOneCredential(ctx, accessToken); err != nil {
		return nil, err
	}

	c.Set("player_id", claims.PlayerID)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}

func (u *middlewareUsecase) Roles(c echo.Context, expectedRoles []int) (echo.Context, error) {
	ctx := c.Request().Context()

	playerRoleCode := c.Get("role_code").(int)

	roleCount, err := u.middlewareRepository.FineRoleCount(ctx)
	if err != nil {
		return nil, err
	}

	playerRole := rbac.IntToBinary(playerRoleCode, int(roleCount))

	for i := 0; i < int(roleCount); i++ {
		if playerRole[i]&expectedRoles[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("error: permission denied")
}
