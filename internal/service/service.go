package service

import (
	"context"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type Storage interface {
	GetAll(ctx context.Context) ([]model.Task, error)
	Create(ctx context.Context, text string) (model.Task, error)
	ToggleDone(ctx context.Context, id int) (model.Task, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (a *Service) List(ctx context.Context) ([]model.Task, error) {
	return a.storage.GetAll(ctx)
}

func (a *Service) Add(ctx context.Context, text string) (model.Task, error) {
	if text == "" {
		return model.Task{}, model.ErrEmptyText
	}
	return a.storage.Create(ctx, text)
}

func (a *Service) Change(ctx context.Context, id int) (model.Task, error) {
	return a.storage.ToggleDone(ctx, id)
}

func (a *Service) Delete(ctx context.Context, id int) error {
	return a.storage.Delete(ctx, id)

}
