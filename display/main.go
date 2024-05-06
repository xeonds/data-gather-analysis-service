package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	router := gin.Default()
	router.GET("/ws", handleWebSocket(conn))
	router.GET("/count", handleGetCount(config))
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("dist/"))))

	panic(router.Run(fmt.Sprint(":", config.Port.Display)))
}

func handleWebSocket(conn *amqp.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		srcConn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer srcConn.Close()

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
			err := srcConn.WriteMessage(websocket.TextMessage, msg.Body)
			if err != nil {
				log.Println(err)
				return
			}
		}
		select {}
	}
}

func handleGetCount(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"count": config.DetectorCount})
	}
}
