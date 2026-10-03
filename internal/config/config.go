package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen Listen
	JWT    JWTConfig
}

type Listen struct {
	BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
	Port   string `env:"PORT" env-default:"8080"`
}

type JWTConfig struct {
	Secret          string        `env:"JWT_SECRET" env-default:"qvH5gTU8zDngVhxVxl2EORoi2Z9P6rgksNfqOiC8pQU"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TTL" env-default:"24h"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"168h"`
}

var instance *Config
var once sync.Once

func LoadConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			fmt.Println("failed loading config", err.Error())
		}
	})

	fmt.Println("config loaded")
	return instance
}
