package middlewareUsecases

import (
	"errors"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/middleware"
	"github.com/korvised/ilog-shop/modules/middleware/middlewareRepositories"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, accessToken string) (echo.Context, error)
		Roles(c echo.Context, expectedRoles []int) error
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

	c.Set(middleware.PlayerID, claims.PlayerID)
	c.Set(middleware.RoleCode, claims.RoleCode)

	return c, nil
}

func (u *middlewareUsecase) Roles(c echo.Context, allowedRoles []int) error {
	playerRoleCode := c.Get(middleware.RoleCode).(int)

	for _, role := range allowedRoles {
		if playerRoleCode == role {
			return nil
		}
	}

	return errors.New("error: permission denied")
}
