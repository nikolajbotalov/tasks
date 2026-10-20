package config

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen   Listen
	JWT      JWT
	Database Database
	CORS     CORS
}

type Listen struct {
	BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
	Port   string `env:"PORT" env-default:"8080"`
}

type JWT struct {
	Secret         string        `env:"JWT_SECRET" env-required:"true"`
	AccessTokenTTL time.Duration `env:"JWT_ACCESS_TTL" env-default:"24h"`
}

type Database struct {
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBName     string `env:"DB_NAME" env-default:"taskflow"`
	DBUser     string `env:"DB_USER" env-default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`
	DBSSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

type CORS struct {
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods string `env:"CORS_ALLOWED_METHODS"`
}

var (
	instance *Config
	once     sync.Once
	initErr  error
)

func LoadConfig() (*Config, error) {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				initErr = fmt.Errorf("read config: %w", err)
				return
			}

			if err = cleanenv.ReadEnv(instance); err != nil {
				initErr = fmt.Errorf("read env: %w", err)
				return
			}
		}

		var errs []error
		if instance.JWT.Secret == "" {
			errs = append(errs, errors.New("JWT_SECRET is empty"))
		}
		if instance.Database.DBPassword == "" {
			errs = append(errs, errors.New("DB_PASSWORD is empty"))
		}
		if instance.CORS.AllowedOrigins == "" {
			errs = append(errs, errors.New("CORS_ALLOWED_ORIGINS is empty"))
		}
		if instance.CORS.AllowedMethods == "" {
			errs = append(errs, errors.New("CORS_ALLOWED_METHODS is empty"))
		}
		initErr = errors.Join(errs...)
	})

	return instance, initErr
}
