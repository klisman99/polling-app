package repositories

import (
	"context"
	"errors"
	"polling-app/poll-service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PollRepository struct {
	collection *mongo.Collection
}

func NewPollRepository(db *mongo.Database) *PollRepository {
	return &PollRepository{
		collection: db.Collection("polls"),
	}
}

func (r *PollRepository) Create(ctx context.Context, poll *models.Poll) error {
	_, err := r.collection.InsertOne(ctx, poll)
	return err
}

func (r *PollRepository) FindAll(ctx context.Context) ([]*models.Poll, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var polls []*models.Poll
	if err := cursor.All(ctx, polls); err != nil {
		return nil, err
	}

	return polls, nil
}

func (r *PollRepository) FindById(ctx context.Context, id string) (*models.Poll, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var poll models.Poll
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&poll)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("poll not found")
		}
		return nil, err
	}

	return &poll, nil
}

func (r *PollRepository) Update(ctx context.Context, id string, poll *models.Poll) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	poll.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"question":   poll.Question,
			"options":    poll.Options,
			"is_active":  poll.IsActive,
			"updated_at": poll.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("poll not found")
	}

	return nil
}

func (r *PollRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("poll not found")
	}

	return nil
}

func (r *PollRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
