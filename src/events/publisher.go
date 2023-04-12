package events

import (
	"encoding/json"
	"fmt"
	"time"

	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/resource"
)

// Define Topic Prefix
const TopicPrefix = "cleovon/test"

// Configuration parameters
var (
	brokerConfig = config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                "tcp://192.168.0.100:55555",
		config.ServicePropertyVPNName:                    "zecomeia",
		config.AuthenticationPropertySchemeBasicUserName: "catatau",
		config.AuthenticationPropertySchemeBasicPassword: "catatau",
	}
	messagingService    solace.MessagingService
	persistentPublisher solace.PersistentMessagePublisher
)

// Receipt Handler
func PublishReceiptListener(receipt solace.PublishReceipt) {
	fmt.Println("Received a Publish Receipt from the broker")
	// fmt.Println("IsPersisted: ", receipt.IsPersisted())
	fmt.Println("Message : ", receipt.GetMessage())
	if receipt.GetError() != nil {
		fmt.Println("Gauranteed Message is NOT persisted on the broker! Received NAK")
		fmt.Println("Error is: ", receipt.GetError())
		// probably want to do something here.  some error handling possibilities:
		//  - send the message again
		//  - send it somewhere else (error handling queue?)
		//  - log and continue
		//  - pause and retry (backoff) - maybe set a flag to slow down the publisher
	}
}

func ConnectPublisher() {
	if mS, err := messaging.NewMessagingServiceBuilder().FromConfigurationProvider(brokerConfig).Build(); err != nil {
		panic(err)
	} else {
		messagingService = mS
	}

	// Connect to the messaging serice
	if err := messagingService.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the broker? ", messagingService.IsConnected())

	//  Build a Persistent Message Publisher
	if p, err := messagingService.CreatePersistentMessagePublisherBuilder().Build(); err != nil {
		panic(err)
	} else {
		persistentPublisher = p
	}

	// Set the message publisher receipt listener
	persistentPublisher.SetMessagePublishReceiptListener(PublishReceiptListener)

	startErr := persistentPublisher.Start()
	if startErr != nil {
		panic(startErr)
	}

	fmt.Println("Persistent Publisher running? ", persistentPublisher.IsRunning())

	fmt.Println("\n===Interrupt (CTR+C) to stop publishing===")
}

func PublishMessage(message interface{}, event string) error {

	if persistentPublisher.IsReady() {
		//  Prepare outbound message payload and body
		messageBody, _ := json.Marshal(message)

		messageBuilder := messagingService.MessageBuilder().
			WithProperty("application", "samples").
			WithProperty("language", "go")
		message, err := messageBuilder.BuildWithStringPayload(string(messageBody))
		if err != nil {
			panic(err)
		}

		topic := resource.TopicOf(TopicPrefix + "/book/" + event)
		fmt.Printf("Publishing on: %s, please ensure queue has matching subscription.\n", topic.GetName())

		// Publish on dynamic topic with dynamic body
		// NOTE: publishing to topic, so make sure GuaranteedSubscriber queue is subscribed to same topic,
		//       or enable "Reject Message to Sender on No Subscription Match" the client-profile
		publishErr := persistentPublisher.Publish(message, topic, nil, nil)
		// Block until message is acknowledged
		// publishErr := persistentPublisher.PublishAwaitAcknowledgement(message, topic, 2*time.Second, nil)

		if publishErr != nil {
			return publishErr
		}

		return nil
	}

	return fmt.Errorf("%s", "Producer not enabled.")
}

func ClosePublisher() {
	// Terminate the Direct Receiver
	persistentPublisher.Terminate(1 * time.Second)
	fmt.Println("\nDirect Publisher Terminated? ", persistentPublisher.IsTerminated())
	// Disconnect the Message Service
	messagingService.Disconnect()
	fmt.Println("Messaging Service Disconnected? ", !messagingService.IsConnected())
}
