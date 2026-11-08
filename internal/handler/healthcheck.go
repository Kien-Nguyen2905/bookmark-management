package handler

import (
	"net/http"

	"bookmark-management/internal/service"

	"github.com/gin-gonic/gin"
)

// HealthCheck is the HTTP handler of the health-check endpoint.
type HealthCheck interface {
	Check(c *gin.Context)
}

type healthCheck struct {
	svc service.HealthCheck
}

func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{svc: svc}
}

// @Summary Health check
// @Description Check the service and its connection to Redis
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthCheckResponse
// @Failure 500 {object} map[string]string
// @Router /health-check [get]
func (h *healthCheck) Check(c *gin.Context) {
	resp, err := h.svc.Check(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
