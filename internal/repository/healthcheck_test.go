package repository

import (
	"testing"

	redisPkg "bookmark-management/internal/pkg/redis"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(c *redis.Client)

		expectedErr bool
	}{
		{
			name: "normal case - redis answers",
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

			mock := redisPkg.InitMockRedis(t)
			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			testRepo := NewHealthCheck(mock)
			err := testRepo.Ping(t.Context())

			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
