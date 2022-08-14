package cache

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string, getAndDelete bool) (string, error)
	Set(ctx context.Context, key string, value string, expiredIn time.Duration) error
}
