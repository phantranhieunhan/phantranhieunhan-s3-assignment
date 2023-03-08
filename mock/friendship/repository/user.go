package mockfriendshiprepo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error) {
	args := m.Called(ctx, emails)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockUserRepository) GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	args := m.Called(ctx, userIDs)
	return args.Get(0).(map[string]string), args.Error(1)
}
