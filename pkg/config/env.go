package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode    string `env:"MODE" envDefault:"development"`
	Port    int    `env:"PORT" envDefault:"8080"`
	GinMode string `env:"GIN_MODE" envDefault:"debug"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
