package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DSN        string
	ServerPort string
	logger     *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *Config {
	return &Config{
		logger: l,
	}
}

func (c *Config) ParseFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		c.logger.Errorf("Error loading.env file: %v", err)
	}

	c.DSN = os.Getenv("DSN")
	c.ServerPort = os.Getenv("SERVER_PORT")

	return c
}
