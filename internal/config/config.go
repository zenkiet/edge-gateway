package config

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Port         string
	Dir          string
	UpstreamURL  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Config struct {
	App AppConfig
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port:         getEnv("APP_PORT", "8080"),
			Dir:          getEnv("DIR", "./dist"),
			UpstreamURL:  getEnv("UPSTREAM_URL", "http://192.168.1.194:5000"),
			ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 10) * time.Second,
			WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 10) * time.Second,
			IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", 120) * time.Second,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return fallback
}
