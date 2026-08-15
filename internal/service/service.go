package service

import (
	"context"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type Storage interface {
	GetAll(ctx context.Context, userID int64) ([]model.Task, error)
	Create(ctx context.Context, userID int64, text string) (model.Task, error)
	ToggleDone(ctx context.Context, userID int64, id int) (model.Task, error)
	Delete(ctx context.Context, userID int64, id int) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (a *Service) List(ctx context.Context, userID int64) ([]model.Task, error) {
	return a.storage.GetAll(ctx, userID)
}

func (a *Service) Add(ctx context.Context, userID int64, text string) (model.Task, error) {
	if text == "" {
		return model.Task{}, model.ErrEmptyText
	}
	return a.storage.Create(ctx, userID, text)
}

func (a *Service) Change(ctx context.Context, userID int64, id int) (model.Task, error) {
	return a.storage.ToggleDone(ctx, userID, id)
}

func (a *Service) Delete(ctx context.Context, userID int64, id int) error {
	return a.storage.Delete(ctx, userID, id)

}
