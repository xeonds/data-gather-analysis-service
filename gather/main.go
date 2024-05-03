package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"data-gather-analysis-service/model"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func detector(id int) func(*amqp.Channel, amqp.Queue) {
	return func(ch *amqp.Channel, q amqp.Queue) {
		msg := model.Data{
			ID:   id,
			Data: rand.Float64(),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		for range time.Tick(time.Second) {
			if err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				},
			); err != nil {
				panic(err)
			}
		}
	}
}

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

	// 根据配置模拟若干个采集终端
	for i := range make([]int, config.DetectorCount) {
		go detector(i)(ch, q)
	}

}
