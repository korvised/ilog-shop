package playerUsecases

import (
	"context"
	"errors"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (u *playerUsecase) GetPlayerCredential(c context.Context, email, password string) (*playerPb.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerCredential(c, email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: GetPlayerCredential failed: %v \n", err)
		return nil, errors.New("error: invalid credentials")
	}

	roleCode := 0
	for _, v := range result.PlayerRoles {
		roleCode += v.RoleCode
	}

	return &playerPb.PlayerProfile{
		Id:        result.ID.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.String(),
		UpdatedAt: result.UpdatedAt.String(),
	}, nil
}
