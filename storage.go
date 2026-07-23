package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileStorage struct {
	fileName string
}

type File struct {
	NextID int    `json:"nextID"`
	Tasks  []Task `json:"tasks"`
}

func (s *FileStorage) load() (File, error) {
	var file File
	data, err := os.ReadFile(s.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return File{NextID: 1}, nil
		}
		return file, err
	}
	err = json.Unmarshal(data, &file)
	return file, err
}

func (s *FileStorage) save(file File) error {
	data, err := json.MarshalIndent(file, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(s.fileName, data, 0644)
	return err
}

func (s *FileStorage) GetAll() ([]Task, error) {
	file, err := s.load()
	if err != nil {
		return []Task{}, fmt.Errorf("GetAll: загрузка задач: %w", err)
	}
	return file.Tasks, nil
}

func (s *FileStorage) Create(text string) (Task, error) {
	file, err := s.load()
	if err != nil {
		return Task{}, fmt.Errorf("Create: загрузка задач: %w", err)
	}
	file.Tasks = addTask(file.Tasks, text, file.NextID)
	file.NextID++
	if err := s.save(file); err != nil {
		return Task{}, fmt.Errorf("Create: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return file.Tasks[len(file.Tasks)-1], nil
}

func (s *FileStorage) ToggleDone(id int) (Task, error) {
	file, err := s.load()
	if err != nil {
		return Task{}, fmt.Errorf("ToggleDone: загрузка задач: %w", err)
	}

	if file.Tasks, err = doneTask(file.Tasks, id); err != nil {
		return Task{}, fmt.Errorf("ToggleDone: изменение состояния задачи: %w", err)
	}

	if err := s.save(file); err != nil {
		return Task{}, fmt.Errorf("ToggleDone: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return file.Tasks[findTaskIndex(file.Tasks, id)], nil
}

func (s *FileStorage) Delete(id int) error {
	file, err := s.load()
	if err != nil {
		return fmt.Errorf("Delete: загрузка задач: %w", err)
	}

	file.Tasks, err = deleteTask(file.Tasks, id)
	if err != nil {
		return fmt.Errorf("Delete: удаление задачи: %w", err)
	}

	if err := s.save(file); err != nil {
		return fmt.Errorf("Delete: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return nil
}
