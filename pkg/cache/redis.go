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
	getFunc := c.client.Get
	if getAndDelete {
		getFunc = c.client.GetDel
	}
	value, err := getFunc(ctx, key).Result()
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

func (c RedisClient) HashSet(ctx context.Context, key string, kv map[string]string) error {
	return c.client.HSet(ctx, key, kv).Err()
}

func (c RedisClient) HashGet(ctx context.Context, key, field string) (string, error) {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}
