package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg repo_mocks

// HealthCheck checks that the storage is reachable.
type HealthCheck interface {
	// Ping returns an error if Redis does not answer.
	Ping(ctx context.Context) error
}

type healthCheck struct {
	c *redis.Client
}

func NewHealthCheck(c *redis.Client) HealthCheck {
	return &healthCheck{c: c}
}

func (h *healthCheck) Ping(ctx context.Context) error {
	return h.c.Ping(ctx).Err()
}
