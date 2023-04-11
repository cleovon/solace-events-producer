package business

import (
	"solace-events-producer/src/events"
	"solace-events-producer/src/models"
)

func SendBookRegisterEvent(book models.Book) map[string]interface{} {
	err := events.PublishMessage(book, "create")
	if err != nil {
		return map[string]interface{}{
			"returnCode":    0,
			"returnMessage": "Success",
		}
	}

	return map[string]interface{}{
		"returnCode":    99,
		"returnMessage": err.Error(),
	}
}
