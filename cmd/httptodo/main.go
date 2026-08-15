package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/AlexKorshun/HTTPtodo/internal/api/httpapi"
	ssogrpc "github.com/AlexKorshun/HTTPtodo/internal/clients/sso/grpc"
	"github.com/AlexKorshun/HTTPtodo/internal/config"
	"github.com/AlexKorshun/HTTPtodo/internal/repository/postgres"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
	"github.com/AlexKorshun/HTTPtodo/web"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ssoClient, err := ssogrpc.New(
		context.Background(),
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)

	if err != nil {
		log.Fatal(err)
	}
	pgStorage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	svc := service.New(pgStorage)

	handler := httpapi.NewHandler(svc, ssoClient, cfg.AppID)

	router := httpapi.NewRouter(handler, cfg.AppSecret, cfg.AppID, web.FS())
	server := httpapi.NewServer(cfg.Port, router)
	if err = server.Run(); err != nil {
		log.Fatal(fmt.Errorf("start server: %w", err))
	}

}
