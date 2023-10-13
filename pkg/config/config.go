package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DSN          string
	ServerPort   string
	TemplatesDir string
	logger       *zap.SugaredLogger
}

func New(l *zap.SugaredLogger) *Config {
	return &Config{
		logger: l,
	}
}

func (c *Config) ParseFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		c.logger.Warnf("Error loading.env file: %v", err)
	}

	c.DSN = os.Getenv("DSN")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.TemplatesDir = os.Getenv("TEMPLATE_DIR")

	return c
}
