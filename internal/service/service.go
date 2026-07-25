package service

import (
	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type Storage interface {
	GetAll() ([]model.Task, error)
	Create(text string) (model.Task, error)
	ToggleDone(id int) (model.Task, error)
	Delete(id int) error
}

type TaskService struct {
	storage Storage
}

func NewTaskService(storage Storage) *TaskService {
	return &TaskService{storage: storage}
}

func (a *TaskService) List() ([]model.Task, error) {
	return a.storage.GetAll()
}

func (a *TaskService) Add(text string) (model.Task, error) {
	if text == "" {
		return model.Task{}, model.ErrEmptyText
	}
	return a.storage.Create(text)
}

func (a *TaskService) Change(id int) (model.Task, error) {
	return a.storage.ToggleDone(id)
}

func (a *TaskService) Delete(id int) error {
	return a.storage.Delete(id)

}
