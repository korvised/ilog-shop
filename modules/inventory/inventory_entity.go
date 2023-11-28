package inventory

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Inventory struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"inventory_id,omitempty"`
		PlayerID string             `bson:"player_id" json:"player_id"`
		ItemID   string             `bson:"item_id" json:"item_id"`
	}
)
