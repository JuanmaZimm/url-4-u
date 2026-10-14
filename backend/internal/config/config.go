package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	BaseURL        string
	DBPath         string
	AppEnv         string
	AllowedOrigins []string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./url4u.db"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	originsStr := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if originsStr != "" {
		origins = strings.Split(originsStr, ",")
	} else {
		origins = []string{"http://localhost:5173", "http://localhost:3000", "*"}
	}

	return &Config{
		Port:           port,
		BaseURL:        baseURL,
		DBPath:         dbPath,
		AppEnv:         appEnv,
		AllowedOrigins: origins,
	}
}
