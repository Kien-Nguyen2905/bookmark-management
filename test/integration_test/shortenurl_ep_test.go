package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bookmark-management/internal/api"
	redisPkg "bookmark-management/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenLink(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		body      string
		setupMock func(c *redis.Client)

		expectedStatus int
		verify         func(t *testing.T, body []byte, c *redis.Client)
	}{
		{
			name:           "normal case - code is stored in redis with the requested exp",
			body:           `{"url":"https://example.com/long/path","exp":3600}`,
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, body []byte, c *redis.Client) {
				var resp map[string]string
				assert.NoError(t, json.Unmarshal(body, &resp))
				assert.Equal(t, "Shorten URL generated successfully!", resp["message"])
				assert.Regexp(t, `^[a-zA-Z0-9]{7}$`, resp["code"])

				url, err := c.Get(t.Context(), resp["code"]).Result()
				assert.NoError(t, err)
				assert.Equal(t, "https://example.com/long/path", url)

				ttl, err := c.TTL(t.Context(), resp["code"]).Result()
				assert.NoError(t, err)
				assert.Equal(t, time.Hour, ttl)
			},
		},
		{
			name:           "fail - invalid input",
			body:           `{"url":"string","exp":3600}`,
			expectedStatus: http.StatusBadRequest,
			verify: func(t *testing.T, body []byte, c *redis.Client) {
				assert.JSONEq(t, `{"error":"Invalid Input"}`, string(body))
			},
		},
		{
			name: "fail - redis is down",
			body: `{"url":"https://example.com","exp":3600}`,
			setupMock: func(c *redis.Client) {
				_ = c.Close()
			},
			expectedStatus: http.StatusInternalServerError,
			verify: func(t *testing.T, body []byte, c *redis.Client) {
				assert.JSONEq(t, `{"error":"Internal Server Error"}`, string(body))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisClient := redisPkg.InitMockRedis(t)
			testAPIEngine := api.NewEngine(&api.Config{AppName: "bookmark_service"}, redisClient)

			// the client is closed only after the engine is built, so the request hits a dead Redis
			if tc.setupMock != nil {
				tc.setupMock(redisClient)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			testAPIEngine.GetApp().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			tc.verify(t, w.Body.Bytes(), redisClient)
		})
	}
}
