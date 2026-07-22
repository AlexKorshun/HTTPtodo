package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const fileName = "todos.json"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func loadList() ([]Task, error) {
	var tasks []Task
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return tasks, err
	}
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveList(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, 0644)
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
		return tasks, fmt.Errorf("такой задачи не существует")
	}

	tasks[i].Done = !tasks[i].Done
	return tasks, nil

}

func deleteTask(tasks []Task, index int) ([]Task, error) {
	i := findTaskIndex(tasks, index)
	if i == -1 {
		return tasks, fmt.Errorf("такой задачи не существует")
	}
	tasks = append(tasks[:i], tasks[i+1:]...)
	return tasks, nil

}

func main() {

	http.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := loadList()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось загрузить задачи")
			return
		}
		respondJSON(w, http.StatusOK, tasks)
	})

	http.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := loadList()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось загрузить задачи")
			return
		}
		type bodyJSON struct {
			Text string `json:"text"`
		}
		body := bodyJSON{}
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "не удалось получить JSON")
			return
		}
		if body.Text == "" {
			respondError(w, http.StatusBadRequest, "указан пустой текст")
			return
		}

		tasks = addTask(tasks, body.Text)

		if err := saveList(tasks); err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось сохранить файл")
			return
		}
		respondJSON(w, http.StatusCreated, tasks[len(tasks)-1])
	})

	http.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := loadList()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось загрузить задачи")
			return
		}
		index, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
			return
		}

		if tasks, err = doneTask(tasks, index); err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		if err := saveList(tasks); err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось сохранить файл")
			return
		}
		respondJSON(w, http.StatusOK, tasks[findTaskIndex(tasks, index)])

	})

	http.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := loadList()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось загрузить задачи")
			return
		}
		index, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
			return
		}
		tasks, err = deleteTask(tasks, index)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		if err := saveList(tasks); err != nil {
			respondError(w, http.StatusInternalServerError, "не удалось сохранить файл")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	fmt.Println("сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
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
