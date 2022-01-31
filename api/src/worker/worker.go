package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gok8r/src/queue"
	"log"
	"strconv"
	"time"
)

func printOnError(err error, msg string) {
	if err != nil {
		log.Println("%s: %s", msg, err)
	}
}

// processIncomingJobs processes incoming deliveries from the
// rabbitMQ declared channel.
func processIncomingJobs(jobPool <-chan amqp.Delivery, channel *amqp.Channel) {
	for job := range jobPool {
		respond(channel, job, "Job started")

		var parsed string = string(job.Body)
		log.Printf("Job received: Sleep for %s seconds", parsed)

		duration, _ := strconv.Atoi(parsed)

		time.Sleep(time.Duration(duration) / 2 * time.Second)
		respond(channel, job, "Job halfway done.")
		time.Sleep(time.Duration(duration) / 2 * time.Second)

		respond(channel, job, "Job complete.")
		err := job.Ack(false)
		if err != nil {
			return
		}
	}
}

func respond(channel *amqp.Channel, job amqp.Delivery, response string) {
	err := channel.Publish(
		"",
		job.CorrelationId,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(response),
		})
	log.Println(response)
	if err != nil {
		log.Println(err)
		return
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
		false,
		false,
		false,
		nil,
	)
	printOnError(err, queue.FailedRegConsumer)

	forever := make(chan bool)

	processIncomingJobs(msgs, ch)

	log.Printf(" [*] Awaiting work. To exit press CTRL+C")
	<-forever
}
