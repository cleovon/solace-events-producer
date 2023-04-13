package controllers

import (
	"solace-events-producer/src/business"
	"solace-events-producer/src/models"

	"github.com/gofiber/fiber/v2"
)

func PostBook(c *fiber.Ctx) error {
	var book models.Book
	c.BodyParser(&book)
	returnBusiness := business.SendBookRegisterEvent(book)

	if returnBusiness["returnCode"] == 0 {
		c.Status(200)
	} else {
		c.Status(500)
	}
	return c.JSON(returnBusiness)
}
