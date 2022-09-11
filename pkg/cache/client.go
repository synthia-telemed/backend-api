package cache

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string, getAndDelete bool) (string, error)
	Set(ctx context.Context, key string, value string, expiredIn time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	MultipleGet(ctx context.Context, keys ...string) ([]string, error)
	MultipleSet(ctx context.Context, kv map[string]string) error
	HashSet(ctx context.Context, key string, kv map[string]string) error
	HashGet(ctx context.Context, key, field string) (string, error)
}
