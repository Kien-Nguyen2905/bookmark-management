package service

import (
	"context"
	"errors"
	"testing"

	"bookmark-management/internal/model"
	repo_mocks "bookmark-management/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		serviceName string
		instanceID  string
		setupMock   func(ctx context.Context) *repo_mocks.HealthCheck

		expectedResponse model.HealthCheckResponse
		expectedErr      error
	}{
		{
			name:        "success - redis answers the ping",
			serviceName: "bookmark_service",
			instanceID:  "instance-1",
			setupMock: func(ctx context.Context) *repo_mocks.HealthCheck {
				repo := &repo_mocks.HealthCheck{}
				repo.On("Ping", ctx).Return(nil)
				return repo
			},
			expectedResponse: model.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "instance-1",
			},
		},
		{
			name:        "fail - redis does not answer the ping",
			serviceName: "bookmark_service",
			instanceID:  "instance-1",
			setupMock: func(ctx context.Context) *repo_mocks.HealthCheck {
				repo := &repo_mocks.HealthCheck{}
				repo.On("Ping", ctx).Return(errors.New("connection refused"))
				return repo
			},
			expectedErr: errors.New("connection refused"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupMock(ctx)

			svc := NewHealthCheck(tc.serviceName, tc.instanceID, repo)

			resp, err := svc.Check(ctx)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, resp)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHealthCheck_GenerateInstanceID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &repo_mocks.HealthCheck{}
	repo.On("Ping", ctx).Return(nil)

	svc := NewHealthCheck("bookmark_service", "", repo)

	first, err := svc.Check(ctx)
	assert.NoError(t, err)
	_, err = uuid.Parse(first.InstanceID)
	assert.NoError(t, err)

	// the UUID is generated once, so it does not change between calls
	second, err := svc.Check(ctx)
	assert.NoError(t, err)
	assert.Equal(t, first.InstanceID, second.InstanceID)
}
