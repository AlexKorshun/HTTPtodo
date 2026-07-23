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

func addTask(tasks []Task, text string, id int) []Task {
	task := Task{id, text, false}
	tasks = append(tasks, task)
	return tasks
}

func doneTask(tasks []Task, id int) ([]Task, error) {
	i := findTaskIndex(tasks, id)
	if i == -1 {
		return tasks, ErrNotFound
	}
	tasks[i].Done = !tasks[i].Done
	return tasks, nil
}

func deleteTask(tasks []Task, id int) ([]Task, error) {
	i := findTaskIndex(tasks, id)
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
