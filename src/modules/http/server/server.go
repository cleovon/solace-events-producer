package server

import (
	conf "solace-events-producer/src/modules/configuration"
	"solace-events-producer/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run() {
	app := fiber.New()
	app.Use(logger.New())
	routes.Register(app)
	app.Listen("0.0.0.0:" + conf.Get("port"))
}
