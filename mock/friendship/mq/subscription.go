package mockmq

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionMQ struct {
	mock.Mock
}

func (m *MockSubscriptionMQ) SubscribeUser(ctx context.Context, ds domain.Subscriptions) error {
	args := m.Called(ctx, ds)
	return args.Error(0)
}
