package notification

import (
	"context"
)

type Client interface {
	Send(ctx context.Context, params SendParams, data map[string]string) error
}

type SendParams struct {
	ID    string
	Title string
	Body  string
}
