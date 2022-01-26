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

// ScheduleWork connects to RabbitMQ and schedules a job detailed by the supplied parameter.
func ScheduleWork(secondsOfWork string) bool {
	pConn := MakeConn()
	const taskPool = "taskPool"

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", pConn.User(), pConn.Pass(), pConn.Host(), pConn.Port())
	conn, err := amqp.Dial(dialUrl)
	if err != nil {
		log.Println("%s: %s", FailedConn, err)
		return false
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("%s: %s", FailedChan, err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		taskPool,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("%s: %s", FailedQueue, err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(secondsOfWork),
		})
	if err != nil {
		log.Println("%s: %s", NotPublished, err)
	}
	log.Printf(" [x] Sent %s\n", secondsOfWork)

	return true
}
