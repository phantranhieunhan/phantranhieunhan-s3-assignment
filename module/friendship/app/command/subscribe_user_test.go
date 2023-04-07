package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/phantranhieunhan/s3-assignment/common"
	mockRepo "github.com/phantranhieunhan/s3-assignment/mock/friendship/repository"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command/payload"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase_SubscribeUser_Handle struct {
	name string
	err  error

	getUserIDsByEmailsData  map[string]string
	getUserIDsByEmailsError error

	getSubscriptionData  domain.SubscriptionStatus
	getSubscriptionError error

	withinTransactionError error

	getFriendshipByUserIDsData  domain.Friendship
	getFriendshipByUserIDsError error

	createData  string
	createError error

	updateError error
}

func TestSubscribeUser_Handle(t *testing.T) {
	t.Parallel()
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockSubscriptionRepo := new(mockRepo.MockSubscriptionRepository)
	mockTransaction := new(mockRepo.MockTransaction)

	repoMock := &RepoMock_TestSubscribeUser_Handle{
		mockSubscriptionRepo: mockSubscriptionRepo,
		mockFriendshipRepo:   mockFriendshipRepo,
		mockUserRepo:         mockUserRepo,
		mockTransaction:      mockTransaction,
	}

	h := NewSubscribeUserHandler(mockFriendshipRepo, mockUserRepo, mockSubscriptionRepo, mockTransaction)

	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	friendshipId := "friendship-id"

	errDB := errors.New("some error from db")

	tcs := []TestCase_SubscribeUser_Handle{
		{
			name: "subscriber a user successfully because user did not connect friend and did not subscribe before",

			err:                         nil,
			getUserIDsByEmailsData:      mapEmails,
			getSubscriptionData:         domain.SubscriptionStatusInvalid,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createData:                  friendshipId,
		},
		{
			name: "subscriber a user successfully because user did not connect friend and unsubscribe before",

			err:                         nil,
			getUserIDsByEmailsData:      mapEmails,
			getSubscriptionData:         domain.SubscriptionStatusUnsubscribed,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
		},
		{
			name: "subscriber a user successfully because user friended but unsubscribe before",

			err:                    nil,
			getUserIDsByEmailsData: mapEmails,
			getSubscriptionData:    domain.SubscriptionStatusUnsubscribed,
			getFriendshipByUserIDsData: domain.Friendship{
				Base:     domain.Base{Id: friendshipId},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusFriended,
			},
		},
		{
			name: "subscriber a user successfully because already subscribe",

			err:                    nil,
			getUserIDsByEmailsData: mapEmails,
			getSubscriptionData:    domain.SubscriptionStatusSubscribed,
		},
		{
			name: "subscriber a user fail because friendship is blocked",

			err:                    common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getUserIDsByEmailsData: mapEmails,
			getSubscriptionData:    domain.SubscriptionStatusUnsubscribed,
			getFriendshipByUserIDsData: domain.Friendship{
				Base:     domain.Base{Id: friendshipId},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusBlocked,
			},
			withinTransactionError: common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "subscriber a user fail because get user id by email fail",

			err:                     common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
			getUserIDsByEmailsError: errDB,
			withinTransactionError:  common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user fail because get subscription fail",

			err:                    common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
			getUserIDsByEmailsData: mapEmails,
			getSubscriptionError:   errDB,
			withinTransactionError: common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user successfully because create fail",

			err:                         common.ErrCannotCreateEntity(domain.Subscription{}.DomainName(), errDB),
			getUserIDsByEmailsData:      mapEmails,
			getSubscriptionData:         domain.SubscriptionStatusInvalid,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createError:                 errDB,
			withinTransactionError:      common.ErrCannotCreateEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name: "subscriber a user fail because update fail",

			err:                         common.ErrCannotUpdateEntity(domain.Subscription{}.DomainName(), errDB),
			getUserIDsByEmailsData:      mapEmails,
			getSubscriptionData:         domain.SubscriptionStatusUnsubscribed,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			updateError:                 errDB,
			withinTransactionError:      common.ErrCannotUpdateEntity(domain.Subscription{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			repoMock.prepare(ctx, t, tc)

			err := h.Handle(ctx, payload.SubscriberUserPayloads{
				payload.SubscriberUserPayload{Requestor: emails[0], Target: emails[1]},
			})
			assert.Equal(t, err, tc.err)
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction, mockSubscriptionRepo)
		})
	}
}

type RepoMock_TestSubscribeUser_Handle struct {
	mockSubscriptionRepo *mockRepo.MockSubscriptionRepository
	mockFriendshipRepo   *mockRepo.MockFriendshipRepository
	mockUserRepo         *mockRepo.MockUserRepository
	mockTransaction      *mockRepo.MockTransaction
}

func (r *RepoMock_TestSubscribeUser_Handle) prepare(ctx context.Context, t *testing.T, tc TestCase_SubscribeUser_Handle) {
	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	subId := "sub-id"

	r.mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()

	if tc.getUserIDsByEmailsError == nil {
		r.mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
			f := args[1].(func(ctx context.Context) error)
			err := f(ctx)
			if tc.withinTransactionError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.withinTransactionError.Error())
			}
		}).Return(tc.withinTransactionError).Once()

		subStatus := tc.getSubscriptionData
		if subStatus.IsNoneExisted() {
			subId = ""
		}
		r.mockSubscriptionRepo.On("GetSubscription", ctx, domain.Subscriptions{
			domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
		}).Return(domain.Subscriptions{
			domain.Subscription{Base: domain.Base{Id: subId}, UserID: friends[1], SubscriberID: friends[0], Status: subStatus},
		}, tc.getSubscriptionError).Once()

		if tc.getSubscriptionError == nil {
			if subStatus.AllowSubscribe() {
				r.mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[1], friends[0]).
					Return(tc.getFriendshipByUserIDsData, tc.getFriendshipByUserIDsError).Once()

				if (tc.getFriendshipByUserIDsError == nil || errors.Is(tc.getFriendshipByUserIDsError, domain.ErrRecordNotFound)) &&
					!tc.getFriendshipByUserIDsData.Status.CanNotSubscribe() {
					if subStatus.IsNoneExisted() {
						r.mockSubscriptionRepo.On("Create", ctx,
							domain.Subscription{UserID: friends[1], SubscriberID: friends[0],
								Status: domain.SubscriptionStatusSubscribed}).
							Return(tc.createData, tc.createError).Once()
					} else {
						r.mockSubscriptionRepo.On("UpdateStatus", ctx, subId, domain.SubscriptionStatusSubscribed).
							Return(tc.updateError).Once()
					}
				}
			}
		}
	}
}
