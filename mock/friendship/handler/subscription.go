package mockHandler

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command/payload"
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

func (m *MockSubscribeUserHandler) Handle(ctx context.Context, payload payload.SubscriberUserPayloads) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

type MockBlockUpdatesUserHandler struct {
	mock.Mock
}

func (m *MockBlockUpdatesUserHandler) Handle(ctx context.Context, payload payload.BlockUpdatesUserPayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

type MockListUpdatesUserHandler struct {
	mock.Mock
}

func (m *MockListUpdatesUserHandler) Handle(ctx context.Context, email string, text string) ([]string, error) {
	args := m.Called(ctx, email, text)
	return args.Get(0).([]string), args.Error(1)
}
