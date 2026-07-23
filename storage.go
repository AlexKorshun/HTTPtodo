package main

import (
	"encoding/json"
	"os"
)

type Storage interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}

type FileStorage struct {
	fileName string
}

func (s *FileStorage) Load() ([]Task, error) {
	var tasks []Task
	data, err := os.ReadFile(s.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return tasks, err
	}
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func (s *FileStorage) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(s.fileName, data, 0644)
	return err
}
