package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookmark-management/internal/model"
	svc_mocks "bookmark-management/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		setupMockSvc func() *svc_mocks.HealthCheck

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success - health check returns OK",
			setupMockSvc: func() *svc_mocks.HealthCheck {
				mockSvc := &svc_mocks.HealthCheck{}
				mockSvc.On("Check", mock.Anything).Return(model.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "bookmark_service",
					InstanceID:  "instance-1",
				}, nil)
				return mockSvc
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-1"}`,
		},
		{
			name: "fail - service returns error",
			setupMockSvc: func() *svc_mocks.HealthCheck {
				mockSvc := &svc_mocks.HealthCheck{}
				mockSvc.On("Check", mock.Anything).Return(model.HealthCheckResponse{}, errors.New("redis is down"))
				return mockSvc
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			mockSvc := tc.setupMockSvc()
			hdl := NewHealthCheck(mockSvc)

			hdl.Check(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			mockSvc.AssertExpectations(t)
		})
	}
}
