package httpserver

import (
	"context"
)

type URLStore interface {
	Save(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, code string) (string, error)
	Delete(ctx context.Context, code string) error
}
