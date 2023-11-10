package authRepositories

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/auth"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepositoryService interface {
		FindCredential(c context.Context, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error)
		FindOneCredential(c context.Context, credentialId string) (*auth.Credential, error)
		FindOneCredentialByAccessToken(c context.Context, accessToken string) (*auth.Credential, error)
		FindOneCredentialByRefreshToken(c context.Context, refreshToken string) (*auth.Credential, error)
		FindOnePlayerProfileToRefresh(c context.Context, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error)
		FindRoleCount(c context.Context) (int64, error)
		InsertCredential(c context.Context, req *auth.Credential) (primitive.ObjectID, error)
		UpdateCredential(c context.Context, req *auth.UpdateCredentialReq) error
		DeleteCredential(c context.Context, credentialId string) error
	}

	authRepository struct {
		db  *mongo.Client
		cfg *config.Config
	}
)

func NewAuthRepository(db *mongo.Client, cfg *config.Config) AuthRepositoryService {
	return &authRepository{db, cfg}
}

func (r *authRepository) authDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}
