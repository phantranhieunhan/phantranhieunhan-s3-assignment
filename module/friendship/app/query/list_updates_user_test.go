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

func TestFriendship_ListUpdatesUserHandler(t *testing.T) {
	t.Parallel()
	emails := []string{"john@example.com", "lisa@example.com", "kate@example.com", "email1@example.com", "email2@example.com", "email3@example.com", "email4@example.com"}
	friends := []string{"friend-1", "friend-2", "friend-3", "friend-4"}

	errDB := errors.New("some error from db")

	tcs := []struct {
		name            string
		result          []string
		requestedEmail  string
		text            string
		mentionedEmails []string

		getUserIDsByEmailsParam string

		getUserIDsByEmailsData  map[string]string
		getUserIDsByEmailsError error

		getSubscriptionSubscribedData  []string
		getSubscriptionSubscribedError error

		getSubscriptionUnsubscribedData  []string
		getSubscriptionUnsubscribedError error

		err error
	}{
		{
			name:                    "get list common friendship successfully",
			text:                    "Hello World!",
			mentionedEmails:         []string{},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getSubscriptionSubscribedData: emails[1:3],
			result:                        emails[1:3],
			err:                           nil,
		},
		{
			name:                    "get list common friendship with valid mention successfully",
			text:                    "Hello World! email1@example.com email2@example.com",
			mentionedEmails:         []string{"email1@example.com", "email2@example.com"},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getSubscriptionSubscribedData:   emails[1:3],
			getSubscriptionUnsubscribedData: emails[5:7],
			result:                          emails[1:5],
			err:                             nil,
		},
		{
			name:                    "get list common friendship fail because invalid mention",
			text:                    "Hello World! email1@example.com email2@example.com",
			mentionedEmails:         []string{"email1@example.com", "email2@example.com"},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsError: domain.ErrNotFoundUserByEmail,
			err:                     common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
		},
		{
			name:                    "get list common friendship fail because getSubscriptionSubscribed has error",
			text:                    "Hello World! email1@example.com email2@example.com",
			mentionedEmails:         []string{"email1@example.com", "email2@example.com"},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getSubscriptionSubscribedError: errDB,
			err:                            common.ErrCannotListEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name:                    "get list common friendship fail because getSubscriptionUnsubscribed has error",
			text:                    "Hello World! email1@example.com email2@example.com",
			mentionedEmails:         []string{"email1@example.com", "email2@example.com"},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getSubscriptionSubscribedData:    emails[1:3],
			getSubscriptionUnsubscribedError: errDB,
			err:                              common.ErrCannotListEntity(domain.Subscription{}.DomainName(), errDB),
		},
		{
			name:                    "get list common friendship successful but mention the blocked user",
			text:                    "Hello World! email1@example.com email2@example.com",
			mentionedEmails:         []string{"email1@example.com", "email2@example.com"},
			getUserIDsByEmailsParam: emails[0],
			getUserIDsByEmailsData: map[string]string{
				emails[0]: friends[0],
				emails[1]: friends[1],
			},
			getSubscriptionSubscribedData:   emails[1:3],
			getSubscriptionUnsubscribedData: emails[4:5],
			result:                          emails[1:4],
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			mockSubscriptionRepo := new(mockRepo.MockSubscriptionRepository)
			mockUserRepo := new(mockRepo.MockUserRepository)
			h := NewListUpdatesUserHandler(mockSubscriptionRepo, mockUserRepo)

			mockUserRepo.On("GetUserIDsByEmails", ctx, []string{tc.getUserIDsByEmailsParam}).Return(tc.getUserIDsByEmailsData, tc.getUserIDsByEmailsError).Once()
			if tc.getUserIDsByEmailsError == nil {
				mockSubscriptionRepo.On("GetSubscriptionEmailsByUserIDAndEmails", ctx, friends[0], tc.mentionedEmails).Return(tc.getSubscriptionSubscribedData, tc.getSubscriptionSubscribedError).Once()
			}
			emails, err := h.Handle(ctx, emails[0], tc.text)
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, emails)

			mock.AssertExpectationsForObjects(t, mockSubscriptionRepo, mockUserRepo)
		})
	}
}
