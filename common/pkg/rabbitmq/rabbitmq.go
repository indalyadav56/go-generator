package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService interface {
	Publish(queueName string, message Message) error
	Consume(queueName string, handler func(message Message) error) error
	CreateQueue(queueName string) error
	Close() error
}

type Message struct {
	Body string
}

type rabbitMQService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(amqpURL string) (RabbitMQService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &rabbitMQService{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *rabbitMQService) Publish(queueName string, message Message) error {
	if _, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Publish the message
	if err := r.channel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
		}); err != nil {
		return err
	}

	log.Println("Message published successfully")
	return nil
}

// func (r *rabbitMQService) Consume(queueName string, handler func(message Message) error) error {

// 	msgs, err := r.channel.Consume(
// 		queueName,
// 		"",    // consumer
// 		false, // auto-ack
// 		false, // exclusive
// 		false, // no-local
// 		false, // no-wait
// 		nil,   // args
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// Process messages in a goroutine
// 	go func() {
// 		for msg := range msgs {
// 			err := handler(Message{Body: string(msg.Body)})
// 			if err != nil {
// 				log.Println("Error handling message:", err)
// 			}
// 		}
// 	}()

// 	return nil
// }

func (r *rabbitMQService) Consume(queueName string, handler func(message Message) error) error {
	// Ensure the queue exists before consuming messages
	if _, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	// Consume messages from the queue
	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages from queue %s: %w", queueName, err)
	}

	// Process messages in a goroutine
	go func() {
		for msg := range msgs {
			err := handler(Message{Body: string(msg.Body)})
			if err != nil {
				log.Println("Error handling message:", err)
			} else {
				// Acknowledge the message if successfully handled
				if err := msg.Ack(false); err != nil {
					log.Println("Failed to acknowledge message:", err)
				}
			}
		}
	}()

	return nil
}

func (r *rabbitMQService) CreateQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue %s: %v", queueName, err)
		return err
	}
	return nil
}

func (r *rabbitMQService) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.connection.Close()
}
