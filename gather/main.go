package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
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

	for {
		data := rand.Float64() // 随机生成数据
		err := ch.Publish(
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
}
