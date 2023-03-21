package query

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

func TestFriendship_ListCommonFriends(t *testing.T) {
	t.Parallel()

	emails := []string{"email-1", "email-2", "email-3", "email-4"}
	friends := []string{"friend-1", "friend-2", "friend-3", "friend-4"}
	requestedEmails := emails[0:2]

	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}

	friendEmails := []string{"email-1", "email-3", "email-2", "email-4", "email-3", "email-4"}
	var nilSlice []string

	errDB := errors.New("some error from db")

	tcs := []struct {
		name            string
		result          []string
		requestedEmails []string
		setup           func(ctx context.Context)

		getUserIDsByEmailsData  map[string]string
		getUserIDsByEmailsError error

		getFriendshipByUserIDAndStatusData  []string
		getFriendshipByUserIDAndStatusError error

		err error
	}{
		{
			name:                               "get list common friendship successfully",
			requestedEmails:                    requestedEmails,
			getUserIDsByEmailsData:             mapEmails,
			getFriendshipByUserIDAndStatusData: friendEmails,
			result:                             emails[2:4],
			err:                                nil,
		},
		{
			name:            "get list common friendship fail because of parameters is invalid",
			requestedEmails: []string{"email"},
			result:          nilSlice,
			err:             common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "emails"),
		},
		{
			name:                    "get list common friendship fail because of get user id by email has error",
			requestedEmails:         requestedEmails,
			getUserIDsByEmailsError: errDB,
			result:                  nilSlice,
			err:                     common.ErrCannotGetEntity(domain.User{}.DomainName(), errDB),
		},
		{
			name:                                "get list common friendship fail because of get friendship by user id and status has error",
			requestedEmails:                     requestedEmails,
			getUserIDsByEmailsData:              mapEmails,
			getFriendshipByUserIDAndStatusError: errDB,
			result:                              nilSlice,
			err:                                 common.ErrCannotListEntity(domain.Friendship{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
			mockUserRepo := new(mockRepo.MockUserRepository)
			h := NewListCommonFriendsHandler(mockFriendshipRepo, mockUserRepo)

			if len(tc.requestedEmails) == EMAIL_TOTAL {
				mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()

				if tc.getUserIDsByEmailsError == nil {
					mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, tc.getUserIDsByEmailsData, []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
						tc.getFriendshipByUserIDAndStatusData, tc.getFriendshipByUserIDAndStatusError).Once()
				}
			}

			ids, err := h.Handle(ctx, tc.requestedEmails)
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, ids)

			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo)
		})
	}
}
