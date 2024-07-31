package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"POSTGRES_NAME" env-required:"true"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type AuthConfig struct {
	User     string `yaml:"user" env:"AUTH_USER" env-required:"true"`
	Password string `yaml:"password" env:"AUTH_PASSWORD" env-required:"true"`
}

type Config struct {
	Environment    string           `yaml:"environment" env:"ENVIRONMENT" env-required:"true"`
	HttpServer     HTTPServerConfig `yaml:"http_server"`
	DatabaseConfig DatabaseConfig   `yaml:"database"`
	StoragePath    string           `yaml:"storage_path" env:"STORAGE_PATH" env-default:"./data/files"`
	AuthConfig     AuthConfig       `yaml:"auth"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH %s doesn't exist", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("failed to read config, err: %v", err)
	}

	return &cfg
}
