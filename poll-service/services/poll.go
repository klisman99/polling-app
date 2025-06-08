package services

import (
	"context"
	"polling-app/poll-service/models"
	"polling-app/poll-service/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PollService struct {
	repo *repositories.PollRepository
}

func NewPollService(repo *repositories.PollRepository) *PollService {
	return &PollService{repo: repo}
}

func (s *PollService) CreatePoll(ctx context.Context, req *models.CreatePollRequest) (*models.Poll, error) {
	options := make([]models.Option, len(req.Options))
	for i, option := range req.Options {
		options[i] = models.Option{
			ID:    primitive.NewObjectID(),
			Text:  option.Text,
			Votes: 0,
		}
	}

	poll := &models.Poll{
		ID:        primitive.NewObjectID(),
		Question:  req.Question,
		Options:   options,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	err := s.repo.Create(ctx, poll)
	if err != nil {
		return nil, err
	}

	return poll, err
}

func (s *PollService) GetAllPolls(ctx context.Context) ([]*models.Poll, error) {
	return s.repo.FindAll(ctx)
}

func (s *PollService) GetPollById(ctx context.Context, id string) (*models.Poll, error) {
	return s.repo.FindById(ctx, id)
}

func (s *PollService) UpdatePoll(ctx context.Context, id string, req *models.UpdatePollRequest) (*models.Poll, error) {
	existingPoll, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	existingPoll.Question = req.Question
	existingPoll.IsActive = req.IsActive

	options := make([]models.Option, len(req.Options))
	for i, option := range req.Options {
		id := option.ID
		if id == primitive.NilObjectID {
			id = primitive.NewObjectID()
		}
		options[i] = models.Option{
			ID:    id,
			Text:  option.Text,
			Votes: option.Votes,
		}
	}

	existingPoll.Options = options

	err = s.repo.Update(ctx, id, existingPoll)
	if err != nil {
		return nil, err
	}

	return existingPoll, nil
}

func (s *PollService) DeletePoll(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
