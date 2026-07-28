package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlexKorshun/HTTPtodo/internal/api"
	"github.com/AlexKorshun/HTTPtodo/internal/config"
	"github.com/AlexKorshun/HTTPtodo/internal/repository/postgres"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	pgStorage, err := postgres.NewPostgresStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewTaskService(pgStorage)
	handler := api.NewHandler(service)

	http.HandleFunc("GET /tasks", handler.GetHandler)
	http.HandleFunc("POST /tasks", handler.PostHandler)
	http.HandleFunc("PATCH /tasks/{id}", handler.PatchHandler)
	http.HandleFunc("DELETE /tasks/{id}", handler.DeleteHandler)

	log.Printf("сервер запущен на :%s\n", cfg.Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))

}
