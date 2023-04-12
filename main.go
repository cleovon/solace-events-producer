package main

import (
	"solace-events-producer/src/events"
	"solace-events-producer/src/modules/http/server"
)

func main() {
	events.ConnectPublisher()
	defer events.ClosePublisher()
	server.Run()
}
