package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/phantranhieunhan/s3-assignment/common"
	mockMq "github.com/phantranhieunhan/s3-assignment/mock/friendship/mq"
	mockRepo "github.com/phantranhieunhan/s3-assignment/mock/friendship/repository"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase_Friendship_ConnectFriendship struct {
	name string
	err  error

	getUserIDsByEmailsError error
	getUserIDsByEmailsData  map[string]string

	withinTransactionError error

	getFriendshipByUserIDsError error
	getFriendshipByUserIDsData  domain.FriendshipStatus

	createError error
	createData  string

	updateError error
}

func TestFriendship_ConnectFriendship(t *testing.T) {
	t.Parallel()
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransaction := new(mockRepo.MockTransaction)
	mockSubMQ := new(mockMq.MockSubscriptionMQ)

	h := NewConnectFriendshipHandler(mockFriendshipRepo, mockUserRepo, mockTransaction, mockSubMQ)

	repoMock := &RepoMock_TestFriendship_ConnectFriendship{
		mockUserRepo:       mockUserRepo,
		mockFriendshipRepo: mockFriendshipRepo,
		mockSubMQ:          mockSubMQ,
		mockTransaction:    mockTransaction,
	}

	emails := []string{"email-1", "email-2"}
	friends := []string{"friend-1", "friend-2"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	friendshipId := "friendship-id"
	friendship := domain.Friendship{UserID: friends[0], FriendID: friends[1]}

	errDB := errors.New("some error from db")

	tcs := []TestCase_Friendship_ConnectFriendship{
		{
			name: "connect friendship successfully because have never connected in the past",

			getUserIDsByEmailsData:      mapEmails,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createData:                  friendshipId,
			err:                         nil,
		},
		{
			name: "connect friendship successfully because they unfriended in the past",

			getUserIDsByEmailsData:     mapEmails,
			getFriendshipByUserIDsData: domain.FriendshipStatusUnfriended,
			err:                        nil,
		},
		{
			name: "connect friendship fail because their relationship is friended",

			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusFriended,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because emails invalid",

			getUserIDsByEmailsError: domain.ErrNotFoundUserByEmail,
			getUserIDsByEmailsData:  make(map[string]string, 0),
			err:                     common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
		},
		{
			name:                       "connect friendship fail because their relationship is blocked",
			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusBlocked,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name:                       "connect friendship fail because their relationship is pending",
			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusPending,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because get friendship by user id fail",

			getUserIDsByEmailsData:      mapEmails,
			withinTransactionError:      common.ErrCannotGetEntity(friendship.DomainName(), errDB),
			getFriendshipByUserIDsError: errDB,
			err:                         common.ErrCannotGetEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because create friendship fail",

			getUserIDsByEmailsData:      mapEmails,
			withinTransactionError:      common.ErrCannotCreateEntity(friendship.DomainName(), errDB),
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createError:                 errDB,
			err:                         common.ErrCannotCreateEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because update friendship fail",

			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrCannotUpdateEntity(friendship.DomainName(), errDB),
			getFriendshipByUserIDsData: domain.FriendshipStatusUnfriended,
			updateError:                errDB,
			err:                        common.ErrCannotUpdateEntity(friendship.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			repoMock.prepare(ctx, t, tc)

			id, err := h.Handle(ctx, emails[0], emails[1])
			assert.Equal(t, err, tc.err)
			if tc.err == nil {
				assert.Equal(t, friendshipId, id)
			}
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction, mockSubMQ)
		})
	}
}

type RepoMock_TestFriendship_ConnectFriendship struct {
	mockUserRepo       *mockRepo.MockUserRepository
	mockFriendshipRepo *mockRepo.MockFriendshipRepository
	mockSubMQ          *mockMq.MockSubscriptionMQ
	mockTransaction    *mockRepo.MockTransaction
}

func (r *RepoMock_TestFriendship_ConnectFriendship) prepare(ctx context.Context, t *testing.T, tc TestCase_Friendship_ConnectFriendship) {
	now := time.Now().UTC()

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
		d := domain.Friendship{
			Base: domain.Base{
				Id:        friendshipId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:   friends[0],
			FriendID: friends[1],
			Status:   tc.getFriendshipByUserIDsData,
		}
		r.mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(d, tc.getFriendshipByUserIDsError).Once()
		if tc.getFriendshipByUserIDsError == domain.ErrRecordNotFound {
			r.mockFriendshipRepo.On("Create", ctx, domain.Friendship{UserID: friends[0], FriendID: friends[1], Status: domain.FriendshipStatusFriended}).Return(friendshipId, tc.createError).Once()
		} else if d.Status.CanConnect() {
			r.mockFriendshipRepo.On("UpdateStatus", ctx, friendshipId, domain.FriendshipStatusFriended).Return(tc.updateError).Once()
		}
		if tc.withinTransactionError == nil {
			r.mockSubMQ.On("SubscribeUser", ctx, domain.Subscriptions{
				domain.Subscription{UserID: friends[0], SubscriberID: friends[1]},
				domain.Subscription{UserID: friends[1], SubscriberID: friends[0]},
			}).Return(nil).Once()
		}
	}
}
