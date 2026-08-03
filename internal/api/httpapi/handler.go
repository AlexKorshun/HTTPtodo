package httpapi

import (
	"context"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type TaskService interface {
	List(ctx context.Context) ([]model.Task, error)
	Add(ctx context.Context, text string) (model.Task, error)
	Change(ctx context.Context, id int) (model.Task, error)
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	taskService TaskService
}

func NewHandler(taskService TaskService) *Handler {
	return &Handler{taskService: taskService}
}
