package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Clients     ClientsConfig
	AppSecret   string
	AppID       int
}

type Client struct {
	Address      string
	Timeout      time.Duration
	RetriesCount int
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

func Load() (Config, error) {
	config := Config{}
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL not set")
	}
	config.DatabaseURL = dbURL
	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT not set")
	}
	config.Port = port

	config.AppSecret = os.Getenv("APP_SECRET")
	if config.AppSecret == "" {
		return Config{}, fmt.Errorf("APP_SECRET not set")
	}

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		return Config{}, fmt.Errorf("APP_ID: %w", err)
	}
	config.AppID = appID

	config.Clients.SSO.Address = os.Getenv("SSO_ADDRESS")
	if config.Clients.SSO.Address == "" {
		return Config{}, fmt.Errorf("SSO_ADDRESS not set")
	}

	timeout, err := time.ParseDuration(os.Getenv("SSO_TIMEOUT"))
	if err != nil {
		return Config{}, fmt.Errorf("SSO_TIMEOUT: %w", err)
	}
	config.Clients.SSO.Timeout = timeout

	retries, err := strconv.Atoi(os.Getenv("SSO_RETRIES"))
	if err != nil {
		return Config{}, fmt.Errorf("SSO_RETRIES: %w", err)
	}
	config.Clients.SSO.RetriesCount = retries

	return config, nil
}
