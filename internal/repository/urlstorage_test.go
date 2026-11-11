package repository

import (
	"testing"
	"time"

	redisPkg "bookmark-management/internal/pkg/redis"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestStoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(c *redis.Client)

		expectedErr bool
	}{
		{
			name: "normal case",
		},
		{
			name: "fail - redis is down",
			setupMock: func(c *redis.Client) {
				_ = c.Close()
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mock := redisPkg.InitMockRedis(t)
			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			testRepo := NewURLStorage(mock)
			err := testRepo.StoreURL(ctx, "abcxyz", "https://google.com", time.Hour)

			if tc.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			res, err := mock.Get(ctx, "abcxyz").Result()
			assert.NoError(t, err)
			assert.Equal(t, "https://google.com", res)

			ttl, err := mock.TTL(ctx, "abcxyz").Result()
			assert.NoError(t, err)
			assert.Equal(t, time.Hour, ttl)
		})
	}
}
