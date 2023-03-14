package friendshiprabbitmq

import (
	"context"
	"encoding/json"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/rabbitmq"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
)

type SubscribeUserMQ struct {
	rabbitMQ rabbitmq.MQ
}

func NewSubscribeUserMQ(rabbitMQ rabbitmq.MQ) SubscribeUserMQ {
	return SubscribeUserMQ{
		rabbitMQ: rabbitMQ,
	}
}

func (s SubscribeUserMQ) SubscribeUser(ctx context.Context, ds domain.Subscriptions) error {
	mByte, err := json.Marshal(ds)
	if err != nil {
		logger.Errorf("[SubscribeUser] Marshal err: %v", err)
		return err
	}

	if err := s.rabbitMQ.PushRawMessage(config.SUBSCRIPTION_EXCHANGE, config.SUBSCRIPTION_CREATED_TOPIC, mByte); err != nil {
		logger.Errorf("[SubscribeUser] PushRawMessage key: %s, data: %v, err: %v", config.SUBSCRIPTION_CREATED_TOPIC, string(mByte), err)
		return err
	}
	return nil
}
