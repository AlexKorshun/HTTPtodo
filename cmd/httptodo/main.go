package main

import (
	"log"
	"net/http"

	"github.com/AlexKorshun/HTTPtodo/internal/api"
	"github.com/AlexKorshun/HTTPtodo/internal/repository/storage"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
)

func main() {

	fileStorage := storage.NewFileStorage("todos.json")
	service := service.NewTaskService(fileStorage)
	handler := api.NewHandler(service)

	http.HandleFunc("GET /tasks", handler.GetHandler)

	http.HandleFunc("POST /tasks", handler.PostHandler)

	http.HandleFunc("PATCH /tasks/{id}", handler.PatchHandler)

	http.HandleFunc("DELETE /tasks/{id}", handler.DeleteHandler)

	log.Println("сервер запущен на :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
