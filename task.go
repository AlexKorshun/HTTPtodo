package main

import "errors"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	ErrEmptyText = errors.New("текст задачи не может быть пустым")
	ErrNotFound  = errors.New("задача не найдена")
)

func addTask(tasks []Task, text string) []Task {
	var task Task
	if len(tasks) == 0 {
		task = Task{1, text, false}
	} else {
		task = Task{tasks[len(tasks)-1].ID + 1, text, false}
	}
	tasks = append(tasks, task)
	return tasks
}

func doneTask(tasks []Task, index int) ([]Task, error) {
	i := findTaskIndex(tasks, index)
	if i == -1 {
		return tasks, ErrNotFound
	}
	tasks[i].Done = !tasks[i].Done
	return tasks, nil
}

func deleteTask(tasks []Task, index int) ([]Task, error) {
	i := findTaskIndex(tasks, index)
	if i == -1 {
		return tasks, ErrNotFound
	}
	tasks = append(tasks[:i], tasks[i+1:]...)
	return tasks, nil

}

func findTaskIndex(tasks []Task, id int) int {
	if id < 0 {
		return -1
	}
	for i, value := range tasks {
		if value.ID == id {
			return i
		}
	}
	return -1
}
