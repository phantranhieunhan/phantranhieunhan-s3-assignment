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
