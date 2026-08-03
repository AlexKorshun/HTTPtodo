package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
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
	return config, nil
}
