package main

import (
	"log"

	"github.com/AlexKorshun/HTTPtodo/internal/api/httpapi"
	"github.com/AlexKorshun/HTTPtodo/internal/config"
	"github.com/AlexKorshun/HTTPtodo/internal/repository/postgres"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pgStorage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	svc := service.New(pgStorage)
	handler := httpapi.NewHandler(svc)
	mux := httpapi.NewRouter(handler)

	server := httpapi.NewServer(cfg.Port, mux)
	log.Fatal(server.Run())

}
