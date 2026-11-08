package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=URLStorage --filename urlstorage.go --outpkg repo_mocks

// URLStorage stores shortened URLs.
type URLStorage interface {
	// StoreURL saves the key -> url pair, which expires after exp.
	StoreURL(ctx context.Context, key, url string, exp time.Duration) error
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

func (s *urlStorage) StoreURL(ctx context.Context, key, url string, exp time.Duration) error {
	return s.c.Set(ctx, key, url, exp).Err()
}
