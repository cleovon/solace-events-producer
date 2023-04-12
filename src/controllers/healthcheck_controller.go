package controllers

import (
	"solace-events-producer/src/business"

	"github.com/gofiber/fiber/v2"
)

func GetHealthCheck(c *fiber.Ctx) error {
	return c.JSON(business.HealthStatus())
}
