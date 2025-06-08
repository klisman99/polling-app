package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Poll struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Question  string             `json:"question" bson:"question"`
	Options   []Option           `json:"options" bson:"options"`
	CreatorID int                `json:"creator_id" bson:"creator_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
}

type Option struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text  string             `json:"text" bson:"text"`
	Votes int                `json:"votes" bson:"votes"`
}

type CreatePollRequest struct {
	Question string   `json:"question" bson:"question"`
	Options  []Option `json:"options" bson:"options"`
}

type UpdatePollRequest struct {
	Question string   `json:"question" bson:"question"`
	Options  []Option `json:"options" bson:"options"`
	IsActive bool     `json:"is_active" bson:"is_active"`
}
