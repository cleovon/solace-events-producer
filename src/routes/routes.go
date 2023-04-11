package routes

import (
	"solace-events-producer/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/healthcheck", controllers.GetHealthCheck)
	app.Post("/book", controllers.PostBook)
}
