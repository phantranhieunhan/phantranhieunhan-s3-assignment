package mockHandler

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/mock"
)

type MockConnectFriendshipHandler struct {
	mock.Mock
}

func (m *MockConnectFriendshipHandler) Handle(ctx context.Context, userEmail, friendEmail string) (domain.Friendship, error) {
	args := m.Called(ctx, userEmail, friendEmail)
	return args.Get(0).(domain.Friendship), args.Error(1)
}

type MockListFriendsHandler struct {
	mock.Mock
}

func (m *MockListFriendsHandler) Handle(ctx context.Context, email string) ([]string, error) {
	args := m.Called(ctx, email)
	return args.Get(0).([]string), args.Error(1)
}

type MockListCommonFriendsHandler struct {
	mock.Mock
}

func (m *MockListCommonFriendsHandler) Handle(ctx context.Context, emails []string) ([]string, error) {
	args := m.Called(ctx, emails)
	return args.Get(0).([]string), args.Error(1)
}
