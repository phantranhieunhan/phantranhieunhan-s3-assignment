package rabbitmq

import (
	"github.com/streadway/amqp"
)

type MQ interface {
	PushKVMessage(exchange, routing string, data KVMessage) error

	PushRawMessage(exchange, routing string, data []byte) error

	Consume() (chan amqp.Delivery, chan error)
}

func New(url string, declarationFile string) (MQ, error) {
	mq := &mqImpl{
		url: url,
	}

	_, err := mq.newConnection()
	if err != nil {
		return nil, err
	}

	err = mq.initFromConfigFile(declarationFile)
	if err != nil {
		return nil, err
	}

	return mq, nil
}
