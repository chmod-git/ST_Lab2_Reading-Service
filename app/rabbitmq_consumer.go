package app

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"testing-project/domain"
)

func startRabbitListener(brokerAddr string) {
	conn, err := amqp.Dial(brokerAddr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"my_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	msgs, err := ch.Consume(
		"my_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %s", err)
	}

	log.Println("Listening for events on RabbitMQ...")

	for msg := range msgs {
		log.Printf("Received: %s", msg.Body)

		var payload struct {
			Event string          `json:"event"`
			Data  *domain.Message `json:"data"`
		}
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			continue
		}

		switch payload.Event {
		case "created", "updated":
			if err := domain.MessageRepo.Save(payload.Data); err != nil {
				log.Printf("Failed to save/update message: %s", err)
			} else {
				log.Printf("Message %s successfully in Redis", payload.Event)
			}
		case "deleted":
			if err := domain.MessageRepo.Delete(payload.Data.Id); err != nil {
				log.Printf("Failed to delete message: %s", err)
			} else {
				log.Printf("Message deleted from Redis")
			}
		default:
			log.Printf("Unknown event type: %s", payload.Event)
		}
	}
}
