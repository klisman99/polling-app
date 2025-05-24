package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Poll struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Question  string             `json:"question" bson:"question"`
	Options   []string           `json:"options" bson:"options"`
	CreatorID int                `json:"creator_id" bson:"creator_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
}
