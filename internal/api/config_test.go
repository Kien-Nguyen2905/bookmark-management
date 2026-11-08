package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The test cases use t.Setenv, so they must not run in parallel.
func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name string

		env map[string]string

		expectedConfig *Config
		expectedErr    bool
	}{
		{
			name:           "default values when env variables are not set",
			expectedConfig: &Config{AppName: "bookmark_service", AppPort: 8080},
		},
		{
			name: "values from env variables",
			env: map[string]string{
				"APP_NAME":    "custom_service",
				"APP_PORT":    "9090",
				"INSTANCE_ID": "instance-1",
			},
			expectedConfig: &Config{AppName: "custom_service", AppPort: 9090, InstanceID: "instance-1"},
		},
		{
			name:        "fail - invalid port",
			env:         map[string]string{"APP_PORT": "abc"},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Setenv restores the original value after the test, os.Unsetenv makes the variable "not set"
			for _, key := range []string{"APP_NAME", "APP_PORT", "INSTANCE_ID"} {
				t.Setenv(key, "")
				assert.NoError(t, os.Unsetenv(key))
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := NewConfig()

			if tc.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedConfig, cfg)
		})
	}
}
