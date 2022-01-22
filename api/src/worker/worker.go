package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gok8r/src/queue"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	pConn := queue.MakeConn()

	const taskPool = "taskPool"

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", pConn.User(), pConn.Pass(), pConn.Host(), pConn.Port())
	conn, err := amqp.Dial(dialUrl)
	failOnError(err, queue.FailedConn)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, queue.FailedChan)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		taskPool,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, queue.FailedQueue)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, queue.FailedRegConsumer)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received: %s", d.Body)
		}
	}()

	log.Printf(" [*] Awaiting work. To exit press CTRL+C")
	<-forever
}
