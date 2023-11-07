package auth

import (
	"github.com/korvised/ilog-shop/modules/player"
	"time"
)

type (
	PlayerLoginReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
	}

	PlayerLogoutReq struct {
		AccessToken string `json:"access_token" form:"access_token" validate:"required,max=500"`
	}

	RefreshTokenReq struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	InsertPlayerRole struct {
		PayerID  string `json:"player_id" validate:"required"`
		RoleCode []int  `json:"role_id" validate:"required"`
	}

	ProfileInterceptor struct {
		*player.PlayerProfile
		Credential *CredentialRes `json:"credential"`
	}

	CredentialRes struct {
		ID           string    `json:"_id"`
		PlayerID     string    `json:"player_id"`
		RoleCode     int       `json:"role_code"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)
