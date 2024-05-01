// data_analyzer.go

package main

import (
	"fmt"

	"github.com/streadway/amqp" // 导入 RabbitMQ 客户端库
)

func main() {
	// 连接 RabbitMQ 服务器
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// 声明一个消息队列
	q, err := ch.QueueDeclare(
		"analysis_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 从数据队列接收数据，执行分析，并发布结果到消息队列
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
		panic(err)
	}

	go func() {
		for msg := range msgs {
			data := string(msg.Body) // 接收到的数据
			// 进行分析计算
			// ...

			// 发布分析结果到消息队列
			err := ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("Analysis Result for %s", data)),
				},
			)
			if err != nil {
				panic(err)
			}
		}
	}()

	// 循环等待
	select {}
}
