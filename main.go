package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	ErrEmptyText = errors.New("текст задачи не может быть пустым")
	ErrNotFound  = errors.New("задача не найдена")
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Storage interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}

type FileStorage struct {
	fileName string
}

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

func main() {

	taskService := TaskService{storage: &FileStorage{fileName: "todos.json"}}
	http.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := taskService.List()
		if err != nil {
			log.Printf("GET /tasks: %v", err)
			respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}
		respondJSON(w, http.StatusOK, tasks)
	})

	http.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		type bodyJSON struct {
			Text string `json:"text"`
		}
		body := bodyJSON{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "не удалось получить JSON")
			return
		}
		task, err := taskService.Add(body.Text)
		if err != nil {
			if errors.Is(err, ErrEmptyText) {
				respondError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Printf("POST /tasks: %v", err)
			respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}
		respondJSON(w, http.StatusCreated, task)
	})

	http.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		index, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
			return
		}

		task, err := taskService.Change(index)

		if err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			log.Printf("PATCH /tasks: %v", err)
			respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}
		respondJSON(w, http.StatusOK, task)

	})

	http.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		index, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
			return
		}

		if err = taskService.Delete(index); err != nil {
			if errors.Is(err, ErrNotFound) {
				respondError(w, http.StatusNotFound, err.Error())
				return
			}
			log.Printf("DELETE /tasks: %v", err)
			respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
