package datastore

import "fmt"

type Config struct {
	Host        string `env:"DATABASE_HOST,required"`
	Port        int    `env:"DATABASE_PORT" envDefault:"5432"`
	User        string `env:"DATABASE_USER,required"`
	Password    string `env:"DATABASE_PASSWORD,required"`
	Name        string `env:"DATABASE_NAME,required"`
	SSLMode     string `env:"DATABASE_SSL_MODE" envDefault:"require"`
	SSLRootCert string `env:"DATABASE_SSL_ROOT_CERT" envDefault:""`
	TimeZone    string `env:"DATABASE_TIMEZONE" envDefault:"Asia/Bangkok"`
}

func (c Config) DSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode, c.TimeZone)
	if len(c.SSLRootCert) != 0 {
		dsn = fmt.Sprintf("%s sslrootcert=%s", dsn, c.SSLRootCert)
	}
	return dsn
}
