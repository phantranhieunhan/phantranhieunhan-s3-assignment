package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	cRabbitMQ "github.com/phantranhieunhan/s3-assignment/common/adapter/rabbitmq"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/middleware"
	"github.com/phantranhieunhan/s3-assignment/module/friendship"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
)

func main() {
	flag.Parse()
	err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
	}

	// Init logger.
	logger.Setup(config.C.Env)
	// Init MQ.
	rabbitmq, err := cRabbitMQ.New(config.C.RabbitMQURL, config.C.MQConfigFile)
	if err != nil {
		logger.Fatal("Failed to connect rabbitmq: ", err)
	}

	db := postgres.NewDatabase()

	r := gin.Default()
	r.Use(middleware.Recover)

	friendship.New(r, db, rabbitmq)

	log.Fatal(r.Run(":" + config.C.Server.Port))
}
