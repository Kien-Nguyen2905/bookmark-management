package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The test cases use t.Setenv, so they must not run in parallel.
func TestNewClient(t *testing.T) {
	testCases := []struct {
		name string

		envPrefix string
		env       map[string]string

		expectedAddr     string
		expectedPassword string
		expectedDB       int
		expectedErr      bool
	}{
		{
			name:         "default values when env variables are not set",
			expectedAddr: "localhost:6379",
		},
		{
			name: "values from env variables",
			env: map[string]string{
				"REDIS_ADDRESS":  "redis:6380",
				"REDIS_PASSWORD": "secret",
				"REDIS_DB":       "2",
			},
			expectedAddr:     "redis:6380",
			expectedPassword: "secret",
			expectedDB:       2,
		},
		{
			name:      "env variables with a prefix",
			envPrefix: "CACHE",
			env: map[string]string{
				"CACHE_REDIS_ADDRESS": "cache:6381",
				"CACHE_REDIS_DB":      "1",
			},
			expectedAddr: "cache:6381",
			expectedDB:   1,
		},
		{
			name:        "fail - invalid db",
			env:         map[string]string{"REDIS_DB": "abc"},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range []string{"REDIS_ADDRESS", "REDIS_PASSWORD", "REDIS_DB"} {
				t.Setenv(key, "")
				assert.NoError(t, os.Unsetenv(key))
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			client, err := NewClient(tc.envPrefix)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, client)
				return
			}

			assert.NoError(t, err)
			defer client.Close()
			assert.Equal(t, tc.expectedAddr, client.Options().Addr)
			assert.Equal(t, tc.expectedPassword, client.Options().Password)
			assert.Equal(t, tc.expectedDB, client.Options().DB)
		})
	}
}
