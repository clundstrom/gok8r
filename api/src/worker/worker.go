package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gok8r/src/queue"
	"log"
	"time"
)

func printOnError(err error, msg string) {
	if err != nil {
		log.Println("%s: %s", msg, err)
	}
}

func main() {
	pConn := queue.MakeConn()

	const taskPool = "taskPool"

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", pConn.User(), pConn.Pass(), pConn.Host(), pConn.Port())
	conn, err := amqp.Dial(dialUrl)
	printOnError(err, queue.FailedConn)
	defer conn.Close()

	ch, err := conn.Channel()
	printOnError(err, queue.FailedChan)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		taskPool,
		false,
		false,
		false,
		false,
		nil,
	)
	printOnError(err, queue.FailedQueue)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	printOnError(err, queue.FailedRegConsumer)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received job: Sleep for %s seconds", d.Body)
			time.Sleep(time.Duration(5) * time.Second)
			log.Printf("Job complete: Sleep for %s seconds", d.Body)
			err := d.Ack(false)
			if err != nil {
				return
			}
		}
	}()

	log.Printf(" [*] Awaiting work. To exit press CTRL+C")
	<-forever
}
