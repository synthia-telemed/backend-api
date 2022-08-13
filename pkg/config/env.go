package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/sms"
)

type Config struct {
	TokenServiceHost string `env:"TOKEN_SERVICE_HOST" envDefault:"localhost:8080"`
	Mode             string `env:"MODE" envDefault:"development"`
	Port             int    `env:"PORT" envDefault:"8080"`
	GinMode          string `env:"GIN_MODE" envDefault:"debug"`
	SentryDSN        string `env:"SENTRY_DSN" envDefault:""`
	DatabaseDSN      string
	DB               datastore.Config
	SMS              sms.Config
	HospitalClient   hospital.Config
	Cache            cache.Config
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
