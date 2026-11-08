package service

import (
	"context"

	"bookmark-management/internal/model"
	"bookmark-management/internal/repository"

	"github.com/google/uuid"
)

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg svc_mocks

// HealthCheck reports the status of the service and its Redis connection.
type HealthCheck interface {
	// Check returns an error if Redis cannot be reached.
	Check(ctx context.Context) (model.HealthCheckResponse, error)
}

type healthCheck struct {
	serviceName string
	instanceID  string
	repo        repository.HealthCheck
}

// NewHealthCheck generates a UUID as the instance id when instanceID is empty.
func NewHealthCheck(serviceName, instanceID string, repo repository.HealthCheck) HealthCheck {
	if instanceID == "" {
		instanceID = uuid.NewString()
	}
	return &healthCheck{serviceName: serviceName, instanceID: instanceID, repo: repo}
}

func (h *healthCheck) Check(ctx context.Context) (model.HealthCheckResponse, error) {
	if err := h.repo.Ping(ctx); err != nil {
		return model.HealthCheckResponse{}, err
	}

	return model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
