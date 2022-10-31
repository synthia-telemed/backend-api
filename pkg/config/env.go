package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/notification"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
)

type Config struct {
	DB             datastore.Config
	SMS            sms.Config
	Payment        payment.Config
	HospitalClient hospital.Config
	GinMode        string `env:"GIN_MODE" envDefault:"debug"`
	SentryDSN      string `env:"SENTRY_DSN" envDefault:""`
	Mode           string `env:"MODE" envDefault:"development"`
	Token          token.Config
	DatabaseDSN    string
	Cache          cache.Config
	Port           int `env:"PORT" envDefault:"8080"`
	Notification   notification.Config
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
