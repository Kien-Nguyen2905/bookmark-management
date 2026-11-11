package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookmark-management/internal/api"
	redisPkg "bookmark-management/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		cfg       *api.Config
		setupMock func(c *redis.Client)

		expectedStatus int
		verifyBody     func(t *testing.T, body []byte)
	}{
		{
			name:           "normal case - service name and instance id from config",
			cfg:            &api.Config{AppName: "bookmark_service", InstanceID: "instance-1"},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body []byte) {
				assert.JSONEq(t, `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-1"}`, string(body))
			},
		},
		{
			name:           "normal case - empty instance id is replaced by a generated uuid",
			cfg:            &api.Config{AppName: "bookmark_service"},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				assert.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "OK", resp["message"])
				_, err := uuid.Parse(resp["instance_id"])
				assert.NoError(t, err)
			},
		},
		{
			name: "fail - redis is down",
			cfg:  &api.Config{AppName: "bookmark_service", InstanceID: "instance-1"},
			setupMock: func(c *redis.Client) {
				_ = c.Close()
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body []byte) {
				assert.JSONEq(t, `{"error":"Internal Server Error"}`, string(body))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisClient := redisPkg.InitMockRedis(t)
			if tc.setupMock != nil {
				tc.setupMock(redisClient)
			}
			testAPIEngine := api.NewEngine(tc.cfg, redisClient)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			testAPIEngine.GetApp().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			tc.verifyBody(t, w.Body.Bytes())
		})
	}
}
