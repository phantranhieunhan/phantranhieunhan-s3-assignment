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

func TestSubscribeUser_Handle(t *testing.T) {
	t.Parallel()
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockSubscriptionRepo := new(mockRepo.MockSubscriptionRepository)
	mockTransaction := new(mockRepo.MockTransaction)

	h := NewSubscribeUserHandler(mockFriendshipRepo, mockUserRepo, mockSubscriptionRepo, mockTransaction)

	// now := time.Now().UTC()

	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	friendshipId := "friendship-id"
	subId := "sub-id"
	// friendship := domain.Friendship{UserID: friends[0], FriendID: friends[1]}

	errDB := errors.New("some error from db")
	var mapNil map[string]string = nil
	var sliceNil []string = nil

	tcs := []struct {
		name  string
		setup func(ctx context.Context)
		err   error
	}{
		{
			name: "subscriber a user successfully because user did not connect friend and did not subscribe before",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusInvalid},
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockSubscriptionRepo.On("Create", ctx, domain.Subscription{UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusSubscribed}).Return(friendshipId, nil).Once()
			},
			err: nil,
		},
		{
			name: "subscriber a user successfully because user did not connect friend and did subscribe and unsubscribe before",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusUnsubscribed},
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockSubscriptionRepo.On("UpdateStatus", ctx, subId, domain.SubscriptionStatusSubscribed).Return(nil).Once()
			},
			err: nil,
		},
		{
			name: "subscriber a user successfully because user friended but unsubscribe before",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusUnsubscribed},
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{
					Base:     domain.Base{Id: friendshipId},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusFriended,
				}, nil).Once()
				mockSubscriptionRepo.On("UpdateStatus", ctx, subId, domain.SubscriptionStatusSubscribed).Return(nil).Once()
			},
			err: nil,
		},
		{
			name: "subscriber a user successfully because already subscribe",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusSubscribed},
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.NoError(t, err)
				}).Return(nil).Once()
			},
			err: nil,
		},
		{
			name: "subscriber a user fail because friendship is blocked",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusUnsubscribed},
				}, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{
					Base:     domain.Base{Id: friendshipId},
					UserID:   friends[0],
					FriendID: friends[1],
					Status:   domain.FriendshipStatusBlocked,
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.Error(t, err)
				}).Return(common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")).Once()
			},
			err: common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "subscriber a user fail because get user id by email fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapNil, sliceNil, errDB).Once()
			},
			err: common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user fail because get subscription fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{}, errDB).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.Error(t, err)
				}).Return(common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB)).Once()
			},
			err: common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user successfully because create fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusInvalid},
				}, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockSubscriptionRepo.On("Create", ctx, domain.Subscription{UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusSubscribed}).Return("", errDB).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.Error(t, err)
				}).Return(common.ErrCannotCreateEntity(domain.Subscription{}.DomainName(), errDB)).Once()
			},
			err: common.ErrCannotCreateEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user fail because update fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(mapEmails, friends, nil).Once()
				mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
					domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
				}).Return(domain.Subscriptions{
					domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusUnsubscribed},
				}, nil).Once()
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
					f := args[1].(func(ctx context.Context) error)
					err := f(ctx)
					assert.Error(t, err)
				}).Return(common.ErrCannotUpdateEntity(domain.Subscription{}.DomainName(), errDB)).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).Return(domain.Friendship{}, domain.ErrRecordNotFound).Once()
				mockSubscriptionRepo.On("UpdateStatus", ctx, subId, domain.SubscriptionStatusSubscribed).Return(errDB).Once()
			},
			err: common.ErrCannotUpdateEntity(domain.Subscription{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			tc.setup(ctx)

			err := h.Handle(ctx, SubscriberUserPayloads{
				SubscriberUserPayload{Requestor: emails[0], Target: emails[1]},
			})
			assert.Equal(t, err, tc.err)
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction, mockSubscriptionRepo)
		})
	}
}
