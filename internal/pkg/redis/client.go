package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a Redis client configured from environment variables.
// A non-empty envPrefix is prepended to the variable names (e.g. CACHE_REDIS_ADDRESS).
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb, nil
}
