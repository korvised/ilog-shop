package playerUsecases

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/modules/player"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/modules/player/playerRepositories"
	"github.com/korvised/ilog-shop/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		AddPlayerMoney(c context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error)
		CreatePlayer(c context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		GetPlayerProfile(c context.Context, playerId string) (*player.PlayerProfile, error)
		GetPlayerSavingAccount(c context.Context, playerId string) (*player.PlayerSavingAccount, error)
		GetPlayerCredential(c context.Context, email, password string) (*playerPb.PlayerProfile, error)
		GetPlayerProfileToRefresh(c context.Context, playerID string) (*playerPb.PlayerProfile, error)
		GetOffset(c context.Context) (int64, error)
		UpsertOffset(c context.Context, offset int64) error
	}

	playerUsecase struct {
		playerRepository playerRepositories.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepositories.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository}
}

func (u *playerUsecase) CreatePlayer(c context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error) {
	if !u.playerRepository.IsUniquePlayer(c, req.Email, req.Username) {
		return nil, errors.New("error: email or username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: failed to hash password")
	}

	// Insert player
	payload := &player.Player{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		PlayerRoles: []player.PlayerRole{
			{RoleTitle: "player", RoleCode: 1},
		},
	}
	playerId, err := u.playerRepository.InsertOnePlayer(c, payload)
	if err != nil {
		return nil, err
	}

	return u.GetPlayerProfile(c, playerId.Hex())
}

func (u *playerUsecase) GetPlayerProfile(c context.Context, playerId string) (*player.PlayerProfile, error) {
	// Find player by id
	result, err := u.playerRepository.FindOnePlayerProfile(c, playerId)
	if err != nil {
		return nil, err
	}

	loc := utils.LoadLocation()

	return &player.PlayerProfile{
		ID:        result.ID.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}

func (u *playerUsecase) AddPlayerMoney(c context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error) {
	// Insert player transaction
	payload := &player.PlayerTransaction{
		PlayerID:  req.PlayerID,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}

	if err := u.playerRepository.InsertOnePlayerTransaction(c, payload); err != nil {
		return nil, err
	}

	return u.GetPlayerSavingAccount(c, req.PlayerID)
}

func (u *playerUsecase) GetPlayerSavingAccount(c context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	return u.playerRepository.FindOnePlayerSavingAccount(c, playerId)
}

func (u *playerUsecase) GetOffset(c context.Context) (int64, error) {
	return u.playerRepository.FindOffset(c)
}

func (u *playerUsecase) UpsertOffset(c context.Context, offset int64) error {
	return u.playerRepository.UpsertOffset(c, offset)
}
