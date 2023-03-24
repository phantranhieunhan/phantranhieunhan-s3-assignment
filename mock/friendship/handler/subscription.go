package mockHandler

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/mock"
)

type MockSubscribeUserHandler struct {
	mock.Mock
}

func (m *MockSubscribeUserHandler) HandleWithSubscription(ctx context.Context, ds domain.Subscriptions) error {
	args := m.Called(ctx, ds)
	return args.Error(0)
}
