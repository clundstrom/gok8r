package queue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const NotPublished = "Could not publish message"
const FailedConn = "Failed to connect to RabbitMQ"
const FailedChan = "Failed to open a channel"
const FailedQueue = "Failed to declare queue"
const FailedRegConsumer = "Failed to register a consumer"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Work(work int) bool {
	pConn := MakeConn()
	const taskPool = "taskPool"

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", pConn.User(), pConn.Pass(), pConn.Host(), pConn.Port())
	conn, err := amqp.Dial(dialUrl)
	failOnError(err, FailedConn)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, FailedChan)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		taskPool,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, FailedQueue)

	body := "Hello World!"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, NotPublished)
	log.Printf(" [x] Sent %s\n", body)

	return true
}
