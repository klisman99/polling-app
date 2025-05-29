package poll

import (
	"context"
	"errors"
	"polling-app/poll-service/internal/db"
	"polling-app/poll-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrEmptyPoll = errors.New("poll question and options cannot be empty")

type Service struct {
	collection *mongo.Collection
}

func NewService(db *db.MongoDB) *Service {
	return &Service{collection: db.Database.Collection("polls")}
}

func (s *Service) Create(ctx context.Context, poll *models.Poll) (string, error) {
	if poll.Question == "" || len(poll.Options) == 0 {
		return "", mongo.ErrNoDocuments
	}

	poll.ID = primitive.NewObjectID()
	poll.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	result, err := s.collection.InsertOne(ctx, poll)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Service) GetAll(ctx context.Context) ([]*models.Poll, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var polls []*models.Poll
	for cursor.Next(ctx) {
		var poll models.Poll
		if err := cursor.Decode(&poll); err != nil {
			return nil, err
		}
		polls = append(polls, &poll)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return polls, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.Poll, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}

	var poll models.Poll
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &poll, nil
}

func (s *Service) Update(ctx context.Context, id string, poll *models.Poll) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	if poll.Question == "" || len(poll.Options) == 0 {
		return ErrEmptyPoll
	}

	poll.ID = objID
	_, err = s.collection.ReplaceOne(ctx, bson.M{"_id": objID}, poll)
	return err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
