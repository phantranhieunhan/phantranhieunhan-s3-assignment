package mockfriendshiprepo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}
