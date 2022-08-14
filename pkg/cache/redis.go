package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Config struct {
	Endpoint  string `env:"REDIS_HOST,required"`
	Username  string `env:"REDIS_USERNAME" envDefault:""`
	Password  string `env:"REDIS_PASSWORD" envDefault:""`
	TLSEnable bool   `env:"REDIS_TLS_ENABLE" envDefault:"false"`
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(config *Config) *RedisClient {
	otp := &redis.Options{
		Addr:     config.Endpoint,
		Username: config.Username,
		Password: config.Password,
	}
	if config.TLSEnable {
		otp.TLSConfig = &tls.Config{}
	}

	return &RedisClient{
		client: redis.NewClient(otp),
	}
}

func (c RedisClient) Get(ctx context.Context, key string, getAndDelete bool) (string, error) {
	var (
		value string
		err   error
	)
	if getAndDelete {
		value, err = c.client.GetDel(ctx, key).Result()
	} else {
		value, err = c.client.Get(ctx, key).Result()
	}
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (c RedisClient) Set(ctx context.Context, key string, value string, expiredIn time.Duration) error {
	return c.client.Set(ctx, key, value, expiredIn).Err()
}
