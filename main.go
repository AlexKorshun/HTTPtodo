package main

import (
	"log"
	"net/http"

	"github.com/AlexKorshun/HTTPtodo/internal/model/repository/storage"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
)

func main() {

	//services := &service.TaskService{storage: &storage.FileStorage{fileName: "todos.json"}}

	fileStorage := storage.NewFileStorage("todos.json")
	service := service.NewTaskService(fileStorage)
	handler := Handler{service}

	http.HandleFunc("GET /tasks", handler.getHandler)

	http.HandleFunc("POST /tasks", handler.postHandler)

	http.HandleFunc("PATCH /tasks/{id}", handler.patchHandler)

	http.HandleFunc("DELETE /tasks/{id}", handler.deleteHandler)

	log.Println("сервер запущен на :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

/*
storage := storage.New()

service := service.New(storage)

handler := handler.New()


*/
