package main_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestConnMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"data_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	data := rand.Float64()
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%f", data)),
		},
	)
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
}
