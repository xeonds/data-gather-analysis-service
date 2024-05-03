package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

func handleWebSocket(conn *amqp.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		srcConn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer srcConn.Close()

		ch, err := conn.Channel()
		if err != nil {
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
			return
		}

		go func() {
			for msg := range msgs {
				err := srcConn.WriteMessage(websocket.TextMessage, msg.Body)
				if err != nil {
					return
				}
			}
		}()
	}
}

func main() {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	router := gin.Default()
	router.GET("/ws", handleWebSocket(conn))
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("dist/"))))

	panic(router.Run(fmt.Sprint(":", config.Port.Display)))
}
