package consumer

import (
	"log"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/rabbitmq"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
	"github.com/streadway/amqp"
)

type HandlerFunc func(data []byte) error

// Consumer : struct
type Consumer struct {
	rabbitMQ rabbitmq.MQ
	app      app.Application
	handlers map[string]HandlerFunc
}

// NewConsumer :
func NewConsumer(rabbitMQ rabbitmq.MQ, application app.Application) *Consumer {
	return &Consumer{
		rabbitMQ: rabbitMQ,
		app:      application,
		handlers: make(map[string]HandlerFunc),
	}
}

func (c *Consumer) Consume() {
	c.setupRoutes()
	msgC, errC := c.rabbitMQ.Consume()
	go func() {
		for {
			select {
			case msg := <-msgC:
				c.process(&msg)
			case err := <-errC:
				log.Println("err", err.Error())
			}
		}
	}()
}

func (c *Consumer) setupRoutes() {
	if c.handlers == nil {
		c.handlers = make(map[string]HandlerFunc)
	}
	c.handlers[config.SUBSCRIPTION_CREATED_TOPIC] = c.SubscribeUserMQ
}

func (c *Consumer) process(msg *amqp.Delivery) {
	defer msg.Ack(false)

	if h, ok := c.handlers[msg.RoutingKey]; ok {
		h(msg.Body)
		return
	}

	logger.Warnf("Consumed unknown message: %s", msg.RoutingKey)
}
