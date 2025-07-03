package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App
		HTTP
		PG
		JWT
		Hasher
	}

	App struct {
		Name    string `env:"APP_NAME"`
		Version string `env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}

	PG struct {
		MaxPoolSize int    `env:"PG_MAX_POOL_SIZE"`
		URL         string `env:"PG_URL"`
	}

	JWT struct {
		SignKey         string        `env:"JWT_SIGN_KEY"`
		AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL"`
	}

	Hasher struct {
		Salt string `env:"HASHER_SALT"`
	}
)

var Cfg *Config

func NewConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("unable to load .env file: %w", err)
	}

	Cfg = &Config{}
	if err := env.Parse(Cfg); err != nil {
		return fmt.Errorf("error parce env: %w", err)
	}

	return nil
}
