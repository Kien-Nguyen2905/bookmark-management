package model

// HealthCheckResponse is the body of GET /health-check.
type HealthCheckResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"18cf5eaf-f401-45bf-a827-265b9c6c9333"`
}
