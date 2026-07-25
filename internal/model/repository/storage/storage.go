package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type FileStorage struct {
	fileName string
}

type File struct {
	NextID int        `json:"nextID"`
	Tasks  []JsonTask `json:"tasks"`
}

func NewFileStorage(fileName string) *FileStorage {
	return &FileStorage{fileName: fileName}
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

func (s *FileStorage) GetAll() ([]model.Task, error) {

	file, err := s.load()
	if err != nil {
		return []model.Task{}, fmt.Errorf("GetAll: загрузка задач: %w", err)
	}

	return convArrayToTask(file.Tasks), nil
}

func (s *FileStorage) Create(text string) (model.Task, error) {
	file, err := s.load()
	if err != nil {
		return model.Task{}, fmt.Errorf("Create: загрузка задач: %w", err)
	}
	file.Tasks = addTask(file.Tasks, text, file.NextID)
	file.NextID++
	if err := s.save(file); err != nil {
		return model.Task{}, fmt.Errorf("Create: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return file.Tasks[len(file.Tasks)-1].convToTask(), nil
}

func (s *FileStorage) ToggleDone(id int) (model.Task, error) {
	file, err := s.load()
	if err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: загрузка задач: %w", err)
	}

	if file.Tasks, err = doneTask(file.Tasks, id); err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: изменение состояния задачи: %w", err)
	}

	if err := s.save(file); err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: сохранение файла: не удалось сохранить файл: %w", err)
	}
	return file.Tasks[findTaskIndex(file.Tasks, id)].convToTask(), nil
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

// БЛЯ надо преобразование из model.Task в jsonTask ну крч похуй сделаешь 😃

func addTask(tasks []JsonTask, text string, id int) []JsonTask {
	task := JsonTask{id, text, false}
	tasks = append(tasks, task)
	return tasks
}

func doneTask(tasks []JsonTask, id int) ([]JsonTask, error) {
	i := findTaskIndex(tasks, id)
	if i == -1 {
		return tasks, model.ErrNotFound
	}
	tasks[i].Done = !tasks[i].Done
	return tasks, nil
}

func deleteTask(tasks []JsonTask, id int) ([]JsonTask, error) {
	i := findTaskIndex(tasks, id)
	if i == -1 {
		return tasks, model.ErrNotFound
	}
	tasks = append(tasks[:i], tasks[i+1:]...)
	return tasks, nil

}

func findTaskIndex(tasks []JsonTask, id int) int {
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
