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
func ScheduleWork(uuid string, secondsOfWork string, messageChan chan string) bool {
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

	// Declare the taskpool queue
	q, err := ch.QueueDeclare(
		taskPool,
		false,
		false,
		false,
		false,
		nil,
	)

	// Response queue unique to user id
	respQ, err := ch.QueueDeclare(
		uuid,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("%s: %s", FailedQueue, err)
	}

	responseQueue, err := ch.Consume(
		respQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	// Publish job and pass response queue via corr id
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: uuid,
			Body:          []byte(secondsOfWork),
		})

	if err != nil {
		log.Println("%s: %s", NotPublished, err)
	}
	log.Printf(" [x] Sent job: %s\n", secondsOfWork)

	for {
		for msg := range responseQueue {
			messageChan <- string(msg.Body)
			log.Println(string(msg.Body))
		}
	}
	return true
}
