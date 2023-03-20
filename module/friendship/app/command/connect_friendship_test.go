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
	mockTransaction := new(mockRepo.MockTransaction)

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

	tcs := []struct {
		name  string
		setup func(ctx context.Context)
		err   error

		getUserIDsByEmailsError error
		getUserIDsByEmailsData  map[string]string

		withinTransactionError error

		getFriendshipByUserIDsError error
		getFriendshipByUserIDsData  domain.FriendshipStatus

		createError error
		createData  string

		updateError error
	}{
		{
			name: "connect friendship successfully because have never connected in the past",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:      mapEmails,
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createData:                  friendshipId,
			err:                         nil,
		},
		{
			name: "connect friendship successfully because they unfriended in the past",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:     mapEmails,
			getFriendshipByUserIDsData: domain.FriendshipStatusUnfriended,
			err:                        nil,
		},
		{
			name: "connect friendship fail because their relationship is friended",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusFriended,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name:                    "connect friendship fail because emails invalid",
			setup:                   func(ctx context.Context) {},
			getUserIDsByEmailsError: domain.ErrNotFoundUserByEmail,
			getUserIDsByEmailsData:  make(map[string]string, 0),
			err:                     common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
		},
		{
			name: "connect friendship fail because their relationship is blocked",
			setup: func(ctx context.Context) {

			},
			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusBlocked,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because their relationship is pending",
			setup: func(ctx context.Context) {

			},
			getUserIDsByEmailsData:     mapEmails,
			withinTransactionError:     common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
			getFriendshipByUserIDsData: domain.FriendshipStatusPending,
			err:                        common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, ""),
		},
		{
			name: "connect friendship fail because get friendship by user id fail",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:      mapEmails,
			withinTransactionError:      common.ErrCannotGetEntity(friendship.DomainName(), errDB),
			getFriendshipByUserIDsError: errDB,
			err:                         common.ErrCannotGetEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because create friendship fail",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:      mapEmails,
			withinTransactionError:      common.ErrCannotCreateEntity(friendship.DomainName(), errDB),
			getFriendshipByUserIDsError: domain.ErrRecordNotFound,
			createError:                 errDB,
			err:                         common.ErrCannotCreateEntity(friendship.DomainName(), errDB),
		},
		{
			name: "connect friendship fail because update friendship fail",
			setup: func(ctx context.Context) {
			},
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
			mockUserRepo := new(mockRepo.MockUserRepository)
			mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)

			mockUserRepo.On("GetUserIDsByEmails", ctx, emails).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()

			if tc.getUserIDsByEmailsError == nil {
				mockTransaction.On("WithinTransaction", ctx, mock.Anything).Run(func(args mock.Arguments) {
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
				mockFriendshipRepo.On("GetFriendshipByUserIDs", ctx, friends[0], friends[1]).Return(d, tc.getFriendshipByUserIDsError).Once()

				if tc.getFriendshipByUserIDsError == domain.ErrRecordNotFound {
					mockFriendshipRepo.On("Create", ctx, domain.Friendship{UserID: friends[0], FriendID: friends[1], Status: domain.FriendshipStatusFriended}).Return(friendshipId, tc.createError).Once()
				} else {
					if d.Status.CanConnect() {
						mockFriendshipRepo.On("UpdateStatus", ctx, friendshipId, domain.FriendshipStatusFriended).Return(tc.updateError).Once()
					}
				}
			}

			tc.setup(ctx)
			h := NewConnectFriendshipHandler(mockFriendshipRepo, mockUserRepo, mockTransaction)
			id, err := h.Handle(ctx, emails[0], emails[1])
			assert.Equal(t, err, tc.err)
			if tc.err == nil {
				assert.Equal(t, friendshipId, id)
			}
			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo, mockTransaction)
		})
	}
}
