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

	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	h := NewListCommonFriendsHandler(mockFriendshipRepo, mockUserRepo)
	emails := []string{"email-1", "email-2", "email-3", "email-4"}
	friends := []string{"friend-1", "friend-2", "friend-3", "friend-4"}
	requestedEmails := emails[0:2]

	mapEmails := map[string]string{
		emails[0]: friends[0],
		emails[1]: friends[1],
	}
	mapUsers := map[string]string{
		friends[1]: emails[1],
		friends[2]: emails[2],
		friends[3]: emails[3],
	}
	var mapNil map[string]string = nil
	var sliceNil []string = nil
	emptySlice := []string{}

	errDB := errors.New("some error from db")

	tcs := []struct {
		name                    string
		result                  []string
		requestedEmails         []string
		setup                   func(ctx context.Context)

		getUserIDsByEmailsData  map[string]string
		getUserIDsByEmailsError error

		err                     error
	}{
		{
			name: "get list common friendship successfully",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(mapEmails, friends[0:2], nil).Once()

				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[0], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{
						{
							UserID:   friends[0],
							FriendID: friends[1],
						},
						{
							UserID:   friends[2],
							FriendID: friends[0],
						},
						{
							UserID:   friends[0],
							FriendID: friends[3],
						},
					}, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[1], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{
						{
							UserID:   friends[0],
							FriendID: friends[1],
						},
						{
							UserID:   friends[2],
							FriendID: friends[1],
						},
						{
							UserID:   friends[1],
							FriendID: friends[3],
						},
					}, nil).Once()
				mockUserRepo.On("GetEmailsByUserIDs", ctx, friends[2:4]).Once().
					Return(mapUsers, emails[2:4], nil)
			},
			requestedEmails: requestedEmails,
			
			result:          emails[2:4],
			err:             nil,
		},
		{
			name: "get list common friendship fail because of parameters is invalid",
			setup: func(ctx context.Context) {
			},
			requestedEmails: []string{"email"},
			result:          emptySlice,
			err:             common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "emails"),
		},
		{
			name: "get list common friendship fail because of get user id by email has error",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(mapNil, sliceNil, errDB).Once()
			},
			requestedEmails: requestedEmails,
			result:          emptySlice,
			err:             common.ErrCannotGetEntity(domain.User{}.DomainName(), errDB),
		},
		{
			name: "get list common friendship fail because of get friendship by user id and status has error",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(mapEmails, friends[0:2], nil).Once()

				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[0], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{}, errDB).Once()
			},
			requestedEmails: requestedEmails,
			result:          emptySlice,
			err:             common.ErrCannotListEntity(domain.Friendship{}.DomainName(), errDB),
		},
		{
			name: "get list common friendship successfully",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(mapEmails, friends[0:2], nil).Once()

				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[0], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{
						{
							UserID:   friends[0],
							FriendID: friends[1],
						},
						{
							UserID:   friends[2],
							FriendID: friends[0],
						},
						{
							UserID:   friends[0],
							FriendID: friends[3],
						},
					}, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[1], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{
						{
							UserID:   friends[0],
							FriendID: friends[1],
						},
						{
							UserID:   friends[2],
							FriendID: friends[1],
						},
						{
							UserID:   friends[1],
							FriendID: friends[3],
						},
					}, nil).Once()
				mockUserRepo.On("GetEmailsByUserIDs", ctx, friends[2:4]).Once().
					Return(mapNil, sliceNil, errDB)
			},
			requestedEmails: requestedEmails,
			result:          emptySlice,
			err:             common.ErrCannotGetEntity(domain.User{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			mockUserRepo.On("GetUserIDsByEmails", ctx, requestedEmails).Return(mapEmails, friends[0:2], nil).Once()

			tc.setup(ctx)

			ids, err := h.Handle(ctx, tc.requestedEmails)
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, ids)

			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo)
		})
	}
}
