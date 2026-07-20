package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
			http.Error(w, "не удалось загрузить задачи", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	http.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := loadList()
		if err != nil {
			http.Error(w, "не удалось загрузить задачи", http.StatusInternalServerError)
			return
		}
		type bodyJSON struct {
			Text string `json:"text"`
		}
		body := bodyJSON{}
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
			http.Error(w, "не удалось получить JSON", http.StatusBadRequest)
			if body.Text == "" {
				http.Error(w, "text пустой", http.StatusBadRequest)
				return
			}
			return
		}
		tasks = addTask(tasks, body.Text)
		fmt.Fprintf(w, "Задача успешно добавлена: %+v \n", tasks[len(tasks)-1])

		if err := saveList(tasks); err != nil {
			http.Error(w, "не удалось сохранить файл", http.StatusInternalServerError)
		}

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
