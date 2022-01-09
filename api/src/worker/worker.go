package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Log if required env not set
func envCheck(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set\n", key)
	}
	return val
}

func main() {

	host := envCheck("HOST")
	port := envCheck("PORT")
	user := envCheck("USER")
	pass := envCheck("PASS")

	const taskPool = "taskPool"

	dialUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	conn, err := amqp.Dial(dialUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		taskPool, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
