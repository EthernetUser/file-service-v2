package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string 
	Port     string 
	User     string 
	Password string 
	Name     string 
}

type HTTPServerConfig struct {
	Address     string        
	Timeout     time.Duration 
	IdleTimeout time.Duration 
}

type AuthConfig struct {
	User     string 
	Password string 
}

type Config struct {
	Environment    string           
	HttpServer     HTTPServerConfig 
	DatabaseConfig DatabaseConfig   
	StoragePath    string           
	AuthConfig     AuthConfig       
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH %s doesn't exist", configPath)
	}

	err := godotenv.Load(configPath)

	if err != nil {
		log.Fatalf("failed to read config, err: %v", err)
	}

	return &Config{
		Environment:    getEnv("ENVIRONMENT", "local"),
		HttpServer: HTTPServerConfig{
			Address:     getEnv("HTTP_SERVER_ADDRESS", "localhost:8080"),
			Timeout:     parseTimeDurationFromEnv("HTTP_SERVER_TIMEOUT", "10s"),
			IdleTimeout: parseTimeDurationFromEnv("HTTP_SERVER_IDLE_TIMEOUT", "120s"),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Name:     getEnv("POSTGRES_NAME", "file-service"),
		},
		StoragePath: getEnv("STORAGE_PATH", "./data/files"),
		AuthConfig: AuthConfig{
			User:     getEnv("AUTH_USER", ""),
			Password: getEnv("AUTH_PASSWORD", ""),
		},
	}
}

func parseTimeDurationFromEnv(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)

	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("failed to parse %s, err: %v", key, err)
	}

	return parsedValue
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue == "" {
		log.Fatalf("environment variable %s not set", key)
	}

	return defaultValue
}