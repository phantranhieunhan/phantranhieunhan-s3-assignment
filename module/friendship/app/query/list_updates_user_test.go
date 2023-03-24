package query

import (
	"context"
	"errors"
	"testing"
	"time"

	mockRepo "github.com/phantranhieunhan/s3-assignment/mock/friendship/repository"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFriendship_ListUpdatesUserHandler(t *testing.T) {
	t.Parallel()
	// text := "Hello World! kate@example.com lisa@example.com lisa@example.com john@example.com"
	emails := []string{"john@example.com", "lisa@example.com", "kate@example.com"}
	friends := []string{"friend-1", "friend-2", "friend-3", "friend-4"}

	// mapEmails := map[string]string{
	// 	emails[0]: friends[0],
	// 	emails[1]: friends[1],
	// }

	friendEmails := []string{"email-1", "email-3", "email-2", "email-4", "email-3", "email-4"}
	var nilSlice []string

	errDB := errors.New("some error from db")

	tcs := []struct {
		name           string
		result         []string
		requestedEmail string
		text           string

		getUserIDsByEmailsData  map[string]string
		getUserIDsByEmailsError error

		getSubscriptionEmailsByUserIDAndStatus  []string
		getSubscriptionEmailsByUserIDAndStatusError error

		err error
	}{
		{
			name: "get list common friendship successfully",
			text: "Hello World! kate@example.com lisa@example.com lisa@example.com",
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getFriendshipByUserIDAndStatusData: friendEmails,
			result:                             emails[2:4],
			err:                                nil,
		},
		// {
		// 	name:   "get list common friendship fail because of parameters is invalid",
		// 	result: nilSlice,
		// 	err:    common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "emails"),
		// },
		// {
		// 	name:                    "get list common friendship fail because of get user id by email has error",
		// 	getUserIDsByEmailsError: errDB,
		// 	result:                  nilSlice,
		// 	err:                     common.ErrCannotGetEntity(domain.User{}.DomainName(), errDB),
		// },
		// {
		// 	name:                                "get list common friendship fail because of get friendship by user id and status has error",
		// 	getUserIDsByEmailsData:              mapEmails,
		// 	getFriendshipByUserIDAndStatusError: errDB,
		// 	result:                              nilSlice,
		// 	err:                                 common.ErrCannotListEntity(domain.Friendship{}.DomainName(), errDB),
		// },
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			mockSubscriptionRepo := new(mockRepo.MockSubscriptionRepository)
			mockUserRepo := new(mockRepo.MockUserRepository)
			h := NewListUpdatesUserHandler(mockSubscriptionRepo, mockUserRepo)

			mockUserRepo.On("GetUserIDsByEmails", ctx, emails[0]).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()
			if tc.getUserIDsByEmailsError == nil {
				mockSubscriptionRepo.On("GetSubscriptionEmailsByUserIDAndStatus", ctx, emails[0], domain.SubscriptionStatusSubscribed).Return(tc.getSubscriptionEmailsByUserIDAndStatus, tc.getSubscriptionEmailsByUserIDAndStatusError)
				
				if tc.getSubscriptionEmailsByUserIDAndStatusError == nil {
                    
                    mockSubscriptionRepo.On("GetSubscriptionEmailsByUserIDAndStatus", ctx, emails[0], domain.SubscriptionStatus
			}
			emails, err := h.Handle(ctx, emails[0], tc.text)
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, emails)

			mock.AssertExpectationsForObjects(t, mockSubscriptionRepo, mockUserRepo)
		})
	}
}
