package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type AuthConfig struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Environment    string           `yaml:"environment"`
	HttpServer     HTTPServerConfig `yaml:"http_server"`
	DatabaseConfig DatabaseConfig   `yaml:"database"`
	StoragePath    string           `yaml:"storage_path"`
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
