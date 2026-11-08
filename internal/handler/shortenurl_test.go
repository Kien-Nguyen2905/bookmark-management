package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	svc_mocks "bookmark-management/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateShortenLink(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		body         string
		setupMockSvc func() *svc_mocks.ShortenUrl

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success - shorten link successfully",
			body: `{"url":"https://example.com","exp":3600}`,
			setupMockSvc: func() *svc_mocks.ShortenUrl {
				mockSvc := &svc_mocks.ShortenUrl{}
				mockSvc.On("CreateShortenLink", mock.Anything, "https://example.com", time.Hour).Return("abc1234", nil)
				return mockSvc
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name:             "fail - url is missing",
			body:             `{"exp":3600}`,
			setupMockSvc:     func() *svc_mocks.ShortenUrl { return &svc_mocks.ShortenUrl{} },
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid Input"}`,
		},
		{
			name:             "fail - url is invalid",
			body:             `{"url":"string","exp":3600}`,
			setupMockSvc:     func() *svc_mocks.ShortenUrl { return &svc_mocks.ShortenUrl{} },
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid Input"}`,
		},
		{
			name:             "fail - exp is missing",
			body:             `{"url":"https://example.com"}`,
			setupMockSvc:     func() *svc_mocks.ShortenUrl { return &svc_mocks.ShortenUrl{} },
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid Input"}`,
		},
		{
			name:             "fail - exp is negative",
			body:             `{"url":"https://example.com","exp":-1}`,
			setupMockSvc:     func() *svc_mocks.ShortenUrl { return &svc_mocks.ShortenUrl{} },
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid Input"}`,
		},
		{
			name:             "fail - body is not valid json",
			body:             `{"url":`,
			setupMockSvc:     func() *svc_mocks.ShortenUrl { return &svc_mocks.ShortenUrl{} },
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid Input"}`,
		},
		{
			name: "fail - service returns error",
			body: `{"url":"https://example.com","exp":3600}`,
			setupMockSvc: func() *svc_mocks.ShortenUrl {
				mockSvc := &svc_mocks.ShortenUrl{}
				mockSvc.On("CreateShortenLink", mock.Anything, "https://example.com", time.Hour).Return("", errors.New("redis is down"))
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
			ginCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(tc.body))
			ginCtx.Request.Header.Set("Content-Type", "application/json")

			mockSvc := tc.setupMockSvc()
			hdl := NewShortenLink(mockSvc)

			hdl.CreateShortenLink(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			mockSvc.AssertExpectations(t)
		})
	}
}
