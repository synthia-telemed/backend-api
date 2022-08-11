package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	TokenServiceHost string `env:"TOKEN_SERVICE_HOST" envDefault:"localhost:8080"`
	Mode             string `env:"MODE" envDefault:"development"`
	Port             int    `env:"PORT" envDefault:"8080"`
	GinMode          string `env:"GIN_MODE" envDefault:"debug"`
	SentryDSN        string `env:"SENTRY_DSN" envDefault:""`
	DatabaseDSN      string
	DB               DatabaseConfig
}

type DatabaseConfig struct {
	Host        string `env:"DATABASE_HOST,required"`
	Port        int    `env:"DATABASE_PORT" envDefault:"5432"`
	User        string `env:"DATABASE_USER,required"`
	Password    string `env:"DATABASE_PASSWORD,required"`
	Name        string `env:"DATABASE_NAME,required"`
	SSLMode     string `env:"DATABASE_SSL_MODE" envDefault:"require"`
	SSLRootCert string `env:"DATABASE_SSL_ROOT_CERT" envDefault:""`
	TimeZone    string `env:"DATABASE_TIMEZONE" envDefault:"Asia/Bangkok"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	cfg.DatabaseDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode, cfg.DB.TimeZone)
	if len(cfg.DB.SSLRootCert) != 0 {
		cfg.DatabaseDSN = fmt.Sprintf("%s sslrootcert=%s", cfg.DatabaseDSN, cfg.DB.SSLRootCert)
	}
	return cfg, err
}
