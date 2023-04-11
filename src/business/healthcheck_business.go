package business

import "solace-events-producer/src/models"

func HealthStatus() models.HealthCheck {
	return models.HealthCheck{
		Status: "The app is healthy",
	}
}
