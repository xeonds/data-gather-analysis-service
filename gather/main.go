package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"data-gather-analysis-service/model"
	"encoding/json"
	"log"
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
			log.Println(err)
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
				log.Println(err)
			}
			log.Println("detector", id, "send data:", msg.Data)
		}
	}
}

func main() {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		panic(err)
	}
	log.Println("connect to rabbitmq success")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	log.Println("open channel success")

	q, err := ch.QueueDeclare(
		"data_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	log.Println("declare queue success")

	// 根据配置模拟若干个采集终端
	log.Println(config.DetectorCount, "detector(s) will be started")
	for i := 0; i < config.DetectorCount; i++ {
		log.Println("start detector", i)
		go detector(i)(ch, q)
	}

	select {}
}
