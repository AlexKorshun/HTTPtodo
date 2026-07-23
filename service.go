package main

import "fmt"

type TaskService struct {
	storage Storage
}

func (a *TaskService) List() ([]Task, error) {
	tasks, err := a.storage.Load()
	if err != nil {
		return []Task{}, fmt.Errorf("list: загрузка задач: %w", err)
	}
	return tasks, nil
}

func (a *TaskService) Add(text string) (Task, error) {
	if text == "" {
		return Task{}, ErrEmptyText
	}
	tasks, err := a.storage.Load()
	if err != nil {
		return Task{}, fmt.Errorf("add: загрузка задач: %w", err)
	}

	tasks = addTask(tasks, text)
	if err := a.storage.Save(tasks); err != nil {
		return Task{}, fmt.Errorf("add: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return tasks[len(tasks)-1], nil
}

func (a *TaskService) Change(index int) (Task, error) {
	tasks, err := a.storage.Load()
	if err != nil {
		return Task{}, fmt.Errorf("change: загрузка задач: %w", err)
	}

	if tasks, err = doneTask(tasks, index); err != nil {
		return Task{}, fmt.Errorf("change: изменение состояния задачи: %w", err)
	}

	if err := a.storage.Save(tasks); err != nil {
		return Task{}, fmt.Errorf("change: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return tasks[findTaskIndex(tasks, index)], nil
}

func (a *TaskService) Delete(index int) error {
	tasks, err := a.storage.Load()
	if err != nil {
		return fmt.Errorf("delete: загрузка задач: %w", err)
	}

	tasks, err = deleteTask(tasks, index)
	if err != nil {
		return fmt.Errorf("delete: удаление задачи: %w", err)
	}

	if err := a.storage.Save(tasks); err != nil {
		return fmt.Errorf("delete: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return nil

}
