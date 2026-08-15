package httpapi

import (
	"context"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type TaskService interface {
	List(ctx context.Context, userID int64) ([]model.Task, error)
	Add(ctx context.Context, userID int64, text string) (model.Task, error)
	Change(ctx context.Context, userID int64, id int) (model.Task, error)
	Delete(ctx context.Context, userID int64, id int) error
}

type Handler struct {
	taskService TaskService
	sso         SSOClient
	appID       int
}

func NewHandler(taskService TaskService, sso SSOClient, appID int) *Handler {
	return &Handler{taskService: taskService, sso: sso, appID: appID}
}
