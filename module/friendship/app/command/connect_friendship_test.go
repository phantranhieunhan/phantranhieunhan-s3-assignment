package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/phantranhieunhan/s3-assignment/common"
	mockRepo "github.com/phantranhieunhan/s3-assignment/mock/friendship/repository"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFriendship_ConnectFriendship(t *testing.T) {
	t.Parallel()
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransaction := new(mockRepo.MockTransaction)
	h := NewConnectFriendshipHandler(mockFriendshipRepo, mockUserRepo, mockTransaction)

	now := time.Now().UTC()

	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	friendshipId := "friendship-id"
	friendship := domain.Friendship{UserID: friends[0], FriendID: friends[1]}

	errDB := errors.New("some error from db")
	var mapNil map[string]string = nil
	var sliceNil []string = nil

	tcs := []struct {
		name  string
		setup func(ctx context.Context)
		err   error
	}{
		{
			name: "connect friendship successfully because have never connected in the past",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
				mockFriendshipRepo.On("Create", ctx, domain.Friendship{UserID: friends[0], FriendID: friends[1], Status: domain.FriendshipStatusFriended}).Return(friendshipId, nil).Once()
			},
			err: nil,
		},
		{
			name: "connect friendship successfully because they unfriended in the past",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{
					Base: domain.Base{
						Id:        friendshipId,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusUnfriended,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
				mockFriendshipRepo.On("UpdateStatus", ctx, friendshipId, domain.FriendshipStatusFriended).Return(nil).Once()
			},
			err: nil,
		},
		{
			name: "connect friendship fail because their relationship is friended",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{
					Base: domain.Base{
						Id:        friendshipId,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusFriended,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Error(t, err.RootError(), domain.ErrFriendshipIsUnavailable)
				}).Return(common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")).Once()
			},
			err: common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because emails invalid",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapNil, sliceNil, domain.ErrNotFoundUserByEmail).Once()
			},
			err: common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
		},
		{
			name: "connect friendship fail because their relationship is blocked",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{
					Base: domain.Base{
						Id:        friendshipId,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusBlocked,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Error(t, err.RootError(), domain.ErrFriendshipIsUnavailable)
				}).Return(common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")).Once()
			},
			err: common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because their relationship is pending",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{
					Base: domain.Base{
						Id:        friendshipId,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusPending,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Error(t, err.RootError(), domain.ErrFriendshipIsUnavailable)
				}).Return(common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")).Once()
			},
			err: common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because get friendship by user id fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{}, errDB).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Equal(t, err.RootError(), errDB)
				}).Return(common.ErrCannotGetEntity(friendship.DomainName(), errDB)).Once()
			},
			err: common.ErrCannotGetEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because create friendship fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Equal(t, err.RootError(), errDB)
				}).Return(common.ErrCannotCreateEntity(friendship.DomainName(), errDB)).Once()
				mockFriendshipRepo.On("Create", ctx, domain.Friendship{UserID: friends[0], FriendID: friends[1], Status: domain.FriendshipStatusFriended}).Return("", common.ErrCannotCreateEntity(friendship.DomainName(), errDB)).Once()
			},
			err: common.ErrCannotCreateEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because update friendship fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(domain.Friendship{
					Base: domain.Base{
						Id:        friendshipId,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusUnfriended,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx).(*common.AppError)
					assert.Equal(t, err.RootError(), errDB)
				}).Return(common.ErrCannotUpdateEntity(friendship.DomainName(), errDB)).Once()
				mockFriendshipRepo.On("UpdateStatus", ctx, friendshipId, domain.FriendshipStatusFriended).Return(common.ErrCannotUpdateEntity(friendship.DomainName(), errDB)).Once()
			},
			err: common.ErrCannotUpdateEntity(friendship.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			tc.setup(ctx)

			id, err := h.Handle(ctx, emails[0], emails[1])
			assert.Equal(t, err, tc.err)
			if tc.err == nil {
				assert.Equal(t, friendshipId, id)
			}
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction)
		})
	}
}
