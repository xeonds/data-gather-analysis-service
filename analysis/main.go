package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"data-gather-analysis-service/model"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	db := lib.NewDB(&config.DatabaseConfig, func(db *gorm.DB) error {
		return db.AutoMigrate(&model.Data{}, &model.Analysis{})
	})
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		log.Println(err)

	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"analysis_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		"data_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	for msg := range msgs {
		recv := new(model.Data)
		if err := json.Unmarshal(msg.Body, recv); err != nil {
			log.Println(err)
		}
		result := analysis(db)(recv)
		data, _ := json.Marshal(result)

		if err := ch.Publish(
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
	}

	select {}
}

func analysis(db *gorm.DB) func(data *model.Data) *model.Analysis {
	return func(data *model.Data) *model.Analysis {
		var sum float64
		db.Create(&data)
		dataSet := make([]model.Data, 0)
		db.Where("id = ?", data.ID).Find(&dataSet)
		for _, d := range dataSet {
			sum += d.Data
		}

		avg := sum / float64(len(dataSet))

		var variance float64
		for _, d := range dataSet {
			variance += (d.Data - avg) * (d.Data - avg)
		}
		variance /= float64(len(dataSet))

		max, min := 0.0, 0.0
		if len(dataSet) > 0 {
			max = dataSet[0].Data
			min = dataSet[0].Data
		}
		for _, d := range dataSet {
			if d.Data > max {
				max = d.Data
			}
			if d.Data < min {
				min = d.Data
			}
		}

		return &model.Analysis{
			ID:       dataSet[0].ID,
			Max:      max,
			Min:      min,
			Avg:      avg,
			Variance: variance,
		}
	}
}
