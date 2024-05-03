package main

import (
	"data-gather-analysis-service/config"
	"data-gather-analysis-service/lib"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	config := lib.LoadConfig[config.Config]()
	conn, err := amqp.Dial(config.MQaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	router := gin.Default()
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("dist/"))))

	panic(router.Run(fmt.Sprint(":", config.Port.Display)))
}
