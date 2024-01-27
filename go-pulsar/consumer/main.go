package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()

	createConsumer(client)
	createReader(client)
}

func createConsumer(client pulsar.Client) {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "topic-1",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	for i := 0; i < 10; i++ {
		// may block here
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
			msg.ID(), string(msg.Payload()))

		consumer.Ack(msg)
	}

	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
}

func createReader(client pulsar.Client) {
	reader, err := client.CreateReader(pulsar.ReaderOptions{
		Topic:          "topic-1",
		StartMessageID: pulsar.EarliestMessageID(),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for {
		if !reader.HasNext() {
			break
		}
		msg, err := reader.Next(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Received message msgId: %#v -- content: '%s'\n", msg.ID(), string(msg.Payload()))
	}
}
