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

type TestCase_Friendship_BlockUpdatesUserHandler struct {
	name string
	err  error

	requestorEmail string
	targetEmail    string

	getUserIDsByEmailsError error
	getUserIDsByEmailsData  map[string]string

	withinTransactionError error

	getFriendshipByUserIDsError error
	getFriendshipByUserIDsData  domain.Friendship

	createError error
	createData  string

	updateError error

	upsertSubscriptionError error
}

func TestFriendship_BlockUpdatesUserHandler(t *testing.T) {
	t.Parallel()
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransaction := new(mockRepo.MockTransaction)
	mockSub := new(mockRepo.MockSubscriptionRepository)

	h := NewBlockUpdatesUserHandler(mockFriendshipRepo, mockUserRepo, mockSub, mockTransaction)

	repoMock := &RepoMock_TestFriendship_BlockUpdatesUserHandler{
		mockUserRepo:         mockUserRepo,
		mockFriendshipRepo:   mockFriendshipRepo,
		mockSubscriptionRepo: mockSub,
		mockTransaction:      mockTransaction,
	}

	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	friendshipId := "friendship-id"

	errDB := errors.New("some error from db")

	tcs := []TestCase_Friendship_BlockUpdatesUserHandler{
		{
			name:           "block updates user successfully because they did not connect friendship before AND did not subscribe before",
			err:            nil,
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData:      mapEmails,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createData:                  friendshipId,
		},
		{
			name:           "block updates user successfully because they did unfriend before AND did not subscribe before",
			err:            nil,
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData: mapEmails,
			getFriendshipByUserIDsData: domain.Friendship{
				Base: domain.Base{
					Id: friendshipId,
				},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusUnfriended,
			},
			createData: friendshipId,
		},
		{
			name:           "block updates user successfully because they did be a friend before AND did not subscribe before",
			err:            nil,
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData: mapEmails,
			getFriendshipByUserIDsData: domain.Friendship{
				Base: domain.Base{
					Id: friendshipId,
				},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusFriended,
			},
		},
		{
			name:           "block updates user fail because they did block together before",
			err:            common.ErrInvalidRequest(domain.ErrCannotBlockUpdatesFromBlockedUser, ""),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData: mapEmails,
			getFriendshipByUserIDsData: domain.Friendship{
				Base: domain.Base{
					Id: friendshipId,
				},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusBlocked,
			},
			withinTransactionError: common.ErrInvalidRequest(domain.ErrCannotBlockUpdatesFromBlockedUser, ""),
		},
		{
			name:           "block updates user fail because GetUserIDsByEmails failed",
			err:            common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsError: errDB,
		},
		{
			name:           "block updates user fail because GetFriendshipByUserIDs failed",
			err:            common.ErrCannotGetEntity(domain.Friendship{}.DomainName(), errDB),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData:      mapEmails,
			getFriendshipByUserIDsError: errDB,
			withinTransactionError:      common.ErrCannotGetEntity(domain.Friendship{}.DomainName(), errDB),
		},
		{
			name:           "block updates user successfully because CREATE failed",
			err:            common.ErrCannotCreateEntity(domain.Friendship{}.DomainName(), errDB),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData:      mapEmails,
			getFriendshipByUserIDsError: errDB,
			createError:                 errDB,
			withinTransactionError:      common.ErrCannotCreateEntity(domain.Friendship{}.DomainName(), errDB),
		},
		{
			name:           "block updates user successfully because Update block friendship failed",
			err:            common.ErrCannotUpdateEntity(domain.Friendship{}.DomainName(), errDB),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData: mapEmails,
			getFriendshipByUserIDsData: domain.Friendship{
				Base: domain.Base{
					Id: friendshipId,
				},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusUnfriended,
			},
			updateError:            errDB,
			withinTransactionError: common.ErrCannotUpdateEntity(domain.Friendship{}.DomainName(), errDB),
		},
		{
			name:           "block updates user successfully because upsert subscription failed",
			err:            common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
			requestorEmail: emails[0],
			targetEmail:    emails[1],

			getUserIDsByEmailsData: mapEmails,
			getFriendshipByUserIDsData: domain.Friendship{
				Base: domain.Base{
					Id: friendshipId,
				},
				UserID:   friends[0],
				FriendID: friends[1],
				Status:   domain.FriendshipStatusUnfriended,
			},
			upsertSubscriptionError: errDB,
			withinTransactionError:  common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			repoMock.prepare(ctx, t, tc)

			err := h.Handle(ctx, BlockUpdatesUserPayload{
				Requestor: tc.requestorEmail,
				Target:    tc.targetEmail,
			})
			assert.Equal(t, err, tc.err)
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction, mockSub)
		})
	}
}

type RepoMock_TestFriendship_BlockUpdatesUserHandler struct {
	mockUserRepo         *mockRepo.MockUserRepository
	mockFriendshipRepo   *mockRepo.MockFriendshipRepository
	mockSubscriptionRepo *mockRepo.MockSubscriptionRepository
	mockTransaction      *mockRepo.MockTransaction
}

func (r *RepoMock_TestFriendship_BlockUpdatesUserHandler) prepare(ctx context.Context, t *testing.T, tc TestCase_Friendship_BlockUpdatesUserHandler) {
	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	friendshipId := "friendship-id"

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
		friendship := tc.getFriendshipByUserIDsData
		r.mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(friendship, tc.getFriendshipByUserIDsError).Once()

		if tc.getFriendshipByUserIDsError == domain.ErrRecordNotFound || friendship.Status.CanBlockUser() {
			r.prepareUnsubscribeMock(ctx, t, tc)
			if tc.upsertSubscriptionError == nil {
				if friendship.Id == "" {
					r.mockFriendshipRepo.On("Create", ctx, domain.Friendship{
						UserID: friends[0], FriendID: friends[1], Status: domain.FriendshipStatusBlocked,
					}).Return(friendshipId, tc.createError).Once()
				} else {
					r.mockFriendshipRepo.On("UpdateStatus", ctx, friendshipId, domain.FriendshipStatusBlocked).Return(tc.updateError).Once()
				}
			}

		} else if friendship.Status == domain.FriendshipStatusFriended {
			r.prepareUnsubscribeMock(ctx, t, tc)
		}
	}

}

func (r *RepoMock_TestFriendship_BlockUpdatesUserHandler) prepareUnsubscribeMock(ctx context.Context, t *testing.T, tc TestCase_Friendship_BlockUpdatesUserHandler) {
	friends := []string{"friend-1", "friend-2"}

	r.mockSubscriptionRepo.On("UnsertSubscription", ctx, domain.Subscription{
		UserID: friends[1], SubscriberID: friends[0], Status: domain.SubscriptionStatusUnsubscribed},
	).Return(tc.upsertSubscriptionError).Once()
}
