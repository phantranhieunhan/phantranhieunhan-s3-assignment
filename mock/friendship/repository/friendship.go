package mockfriendshiprepo

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/mock"
)

type MockFriendshipRepository struct {
	mock.Mock
}

func (m *MockFriendshipRepository) Create(ctx context.Context, d domain.Friendship) (string, error) {
	args := m.Called(ctx, d)
	return args.String(0), args.Error(1)
}

func (m *MockFriendshipRepository) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockFriendshipRepository) GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error) {
	args := m.Called(ctx, userID, friendID)
	return args.Get(0).(domain.Friendship), args.Error(1)
}

func (m *MockFriendshipRepository) GetFriendshipByUserIDAndStatus(ctx context.Context, mapEmailUser map[string]string, status ...domain.FriendshipStatus) ([]string, error) {
	args := m.Called(ctx, mapEmailUser, status)
	return args.Get(0).([]string), args.Error(1)
}
