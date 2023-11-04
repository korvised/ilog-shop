package authRepositories

import (
	"context"
	"errors"
	"github.com/korvised/ilog-shop/modules/auth"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type (
	AuthRepositoryService interface {
		CredentialSearch(c context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error)
		InsertOnePlayerCredential(c context.Context, req *auth.Credential) (primitive.ObjectID, error)
		FindOnePlayerCredential(c context.Context, credentialID string) (*auth.Credential, error)
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db}
}

func (r *authRepository) authDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("auth-db")
}

func (r *authRepository) InsertOnePlayerCredential(c context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerCredential failed: %v \n", err)
		return primitive.NilObjectID, errors.New("error: insert one player credential failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *authRepository) FindOnePlayerCredential(c context.Context, credentialID string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(c, bson.M{"_id": utils.ConvertToObjectId(credentialID)}).Decode(credential); err != nil {
		log.Printf("Error: FindOnePlayerCredential failed: %v \n", err)
		return nil, errors.New("error: find one player credential failed")
	}

	return credential, nil
}
