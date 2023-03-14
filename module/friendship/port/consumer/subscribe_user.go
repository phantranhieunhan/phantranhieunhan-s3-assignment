package consumer

import (
	"context"
	"encoding/json"

	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

func (s *Consumer) SubscribeUserMQ(data []byte) error {
	ctx := context.Background()
	logger.Infof("[Consume] SubscribeUserMQ: %s", string(data))

	var msg domain.Subscriptions
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Error("[Consume] SubsribeUserMQ: unmarshal message failed, ", err)
		return err
	}
	if err := s.app.Commands.SubscribeUser.HandleWithSubscription(ctx, msg); err != nil {
		logger.Errorf("Create Subscription fail when create connection friendship HandleWithSubscription: %w", err)
		return err
	}
	return nil
}
