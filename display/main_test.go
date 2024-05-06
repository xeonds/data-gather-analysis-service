package main_test

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"log"
	"testing"

	"github.com/streadway/amqp"
)

func TestPrintData(t *testing.T) {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"analysis_queue", // queue name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Println(err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err)
		return
	}

	for msg := range msgs {
		log.Println(string(msg.Body))
	}
	select {}
}
