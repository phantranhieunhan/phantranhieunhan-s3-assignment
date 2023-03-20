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

func TestFriendship_ListFriends(t *testing.T) {
	t.Parallel()

	emails := []string{"email-1", "email-2", "email-3", "email-4"}
	friends := []string{"friend-1", "friend-2", "friend-3", "friend-4"}
	mapEmails := map[string]string{
		emails[0]: friends[0],
	}
	mapUsers := map[string]string{
		friends[1]: emails[1],
		friends[2]: emails[2],
		friends[3]: emails[3],
	}
	var mapNil map[string]string = nil
	var sliceNil []string = nil

	errDB := errors.New("some error from db")

	tcs := []struct {
		name                    string
		result                  []string
		setup                   func(ctx context.Context)
		err                     error
		getUserIDsByEmailsError error
		getUserIDsByEmailsData  map[string]string

		getFriendshipByUserIDAndStatusError error
		getFriendshipByUserIDAndStatusData  []string
	}{
		{
			name: "get list friendship successfully",
			setup: func(ctx context.Context) {
			},
			getUserIDsByEmailsData:             mapEmails,
			getFriendshipByUserIDAndStatusData: emails[1:4],
			result:                             emails[1:4],
			err:                                nil,
		},
		{
			name:                    "get list friendship fail because email invalid",
			result:                  nil,
			err:                     common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
			getUserIDsByEmailsError: domain.ErrNotFoundUserByEmail,
			getUserIDsByEmailsData:  map[string]string{},
		},
		{
			name:                                "get list friendship fail because get friendship by userId and status fail",
			getUserIDsByEmailsData:              mapEmails,
			getFriendshipByUserIDAndStatusError: errDB,
			getFriendshipByUserIDAndStatusData:  []string{},
			result:                              nil,
			err:                                 common.ErrCannotListEntity(domain.Friendship{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
			mockUserRepo := new(mockRepo.MockUserRepository)
			h := NewListFriendsHandler(mockFriendshipRepo, mockUserRepo)

			mockUserRepo.On("GetUserIDsByEmails", ctx, []string{emails[0]}).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()
			if tc.getUserIDsByEmailsError == nil {
				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[0:1], emails[0:1], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					tc.getFriendshipByUserIDAndStatusData, tc.getFriendshipByUserIDAndStatusError).Once()
			}

			ids, err := h.Handle(ctx, emails[0])
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, ids)

			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo)
		})
	}
}
