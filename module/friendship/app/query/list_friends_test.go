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
	mockFriendshipRepo := new(mockRepo.MockFriendshipRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)
	h := NewListFriendsHandler(mockFriendshipRepo, mockUserRepo)

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

	errDB := errors.New("some error from db")

	tcs := []struct {
		name   string
		result []string
		setup  func(ctx context.Context)
		err    error
	}{
		{
			name: "get list friendship successfully",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, []string{emails[0]}).Return(mapEmails, nil).Once()
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
				mockUserRepo.On("GetEmailsByUserIDs", ctx, []string{friends[1], friends[2], friends[3]}).Once().
					Return(mapUsers, nil)
			},
			result: []string{emails[1], emails[2], emails[3]},
			err:    nil,
		},
		{
			name: "get list friendship fail because email invalid",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, []string{emails[0]}).Return(map[string]string{}, domain.ErrNotFoundUserByEmail).Once()
			},
			result: nil,
			err:    common.ErrInvalidRequest(domain.ErrNotFoundUserByEmail, "emails"),
		},
		{
			name: "get list friendship fail because get friendship by userId and status fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, []string{emails[0]}).Return(mapEmails, nil).Once()
				mockFriendshipRepo.On("GetFriendshipByUserIDAndStatus", ctx, friends[0], []domain.FriendshipStatus{domain.FriendshipStatusFriended}).Return(
					domain.Friendships{}, errDB).Once()
			},
			result: nil,
			err:    common.ErrCannotListEntity(domain.Friendship{}.DomainName(), errDB),
		},
		{
			name: "get list friendship fail because get emails by userIds fail",
			setup: func(ctx context.Context) {
				mockUserRepo.On("GetUserIDsByEmails", ctx, []string{emails[0]}).Return(mapEmails, nil).Once()
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
					}, nil).Once()
				mockUserRepo.On("GetEmailsByUserIDs", ctx, []string{friends[1], friends[2]}).Once().
					Return(map[string]string{}, errDB)
			},
			result: nil,
			err:    common.ErrCannotGetEntity(domain.User{}.DomainName(), errDB),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			tc.setup(ctx)

			ids, err := h.Handle(ctx, emails[0])
			assert.Equal(t, err, tc.err)
			assert.Equal(t, tc.result, ids)

			mock.AssertExpectationsForObjects(t, mockFriendshipRepo, mockUserRepo)
		})
	}
}
