package mockfriendshiprepo

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, d domain.Subscription) (string, error) {
	args := m.Called(ctx, d)
	return args.String(0), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetSubscription(ctx context.Context, ss domain.Subscriptions) (domain.Subscriptions, error) {
	args := m.Called(ctx, ss)
	return args.Get(0).(domain.Subscriptions), args.Error(1)
}

func (m *MockSubscriptionRepository) UnsertSubscription(ctx context.Context, d domain.Subscription) (string, error) {
	args := m.Called(ctx, d)
	return args.String(0), args.Error(1)
}

func (m *MockSubscriptionRepository) GetSubscriptionEmailsByUserIDAndStatus(ctx context.Context, id string, status domain.SubscriptionStatus) ([]string, error) {
	args := m.Called(ctx, id, status)
	return args.Get(0).([]string), args.Error(1)
}
