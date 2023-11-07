package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	Credential struct {
		ID           primitive.ObjectID `bson:"_id,omitempty"`
		PlayerID     string             `bson:"player_id"`
		RoleCode     int                `bson:"role_code"`
		AccessToken  string             `bson:"access_token"`
		RefreshToken string             `bson:"refresh_token"`
		CreatedAt    time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
	}

	Role struct {
		ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Title string             `json:"title" bson:"title"`
		Code  int                `json:"code" bson:"code"`
	}

	UpdateCredentialReq struct {
		ID           string    `json:"_id" bson:"_id,omitempty"`
		AccessToken  string    `bson:"access_token"`
		RefreshToken string    `bson:"refresh_token"`
		UpdatedAt    time.Time `bson:"updated_at"`
	}
)
