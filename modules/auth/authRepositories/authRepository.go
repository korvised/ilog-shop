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
		InsertCredential(c context.Context, req *auth.Credential) (primitive.ObjectID, error)
		FindOneCredential(c context.Context, credentialId string) (*auth.Credential, error)
		FindOneCredentialByAccessToken(c context.Context, credentialId string) (*auth.Credential, error)
		FindOneCredentialByRefreshToken(c context.Context, credentialId string) (*auth.Credential, error)
		FindOnePlayerProfileToRefresh(c context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error)
		UpdateCredential(c context.Context, req *auth.UpdateCredentialReq) error
		DeleteCredential(c context.Context, credentialId string) error
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db}
}

func (r *authRepository) authDbConn(_ context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *authRepository) InsertCredential(c context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertCredential failed: %v \n", err)
		return primitive.NilObjectID, errors.New("error: insert one player credential failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *authRepository) FindOneCredential(c context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(c, bson.M{"_id": utils.ConvertToObjectId(credentialId)}).Decode(credential); err != nil {
		log.Printf("Error: FindOneCredential failed: %v \n", err)
		return nil, errors.New("error: find one player credential failed")
	}

	return credential, nil
}

func (r *authRepository) FindOneCredentialByAccessToken(c context.Context, accessToken string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(c, bson.M{"access_token": accessToken}).Decode(credential); err != nil {
		log.Printf("Error: FindOneCredentialByAccessToken: %v \n", err)
		return nil, errors.New("error: invalid access token")
	}

	return credential, nil
}

func (r *authRepository) FindOneCredentialByRefreshToken(c context.Context, refreshToken string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(c, bson.M{"refresh_token": refreshToken}).Decode(credential); err != nil {
		log.Printf("Error: FindOneCredentialByRefreshToken: %v \n", err)
		return nil, errors.New("error: invalid refresh token")
	}

	return credential, nil
}

func (r *authRepository) UpdateCredential(c context.Context, req *auth.UpdateCredentialReq) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	if _, err := col.UpdateOne(
		c,
		bson.M{"_id": utils.ConvertToObjectId(req.ID)},
		bson.M{
			"$set": bson.M{
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    req.UpdatedAt,
			},
		}); err != nil {
		log.Printf("Error: UpdateCredential failed: %v \n", err)
		return errors.New("error: update one player credential failed")
	}

	return nil
}

func (r *authRepository) DeleteCredential(c context.Context, credentialId string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	if _, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)}); err != nil {
		log.Printf("Error: DeleteCredential failed: %v \n", err)
		return errors.New("error: delete credential failed")
	}

	return nil
}
