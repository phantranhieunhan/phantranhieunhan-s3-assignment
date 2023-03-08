package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListFriendsRepo interface {
	GetFriendshipByUserIDAndStatus(ctx context.Context, userID string, status ...domain.FriendshipStatus) (domain.Friendships, error)
}

type UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}

type ListFriendsHandler struct {
	repo     ListFriendsRepo
	userRepo UserRepo
}

func NewListFriendsHandler(repo ListFriendsRepo, userRepo UserRepo) ListFriendsHandler {
	return ListFriendsHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h ListFriendsHandler) Handle(ctx context.Context, email string) ([]string, error) {
	var zList []string

	// get userId from email to check available
	userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, []string{email})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return zList, common.ErrInvalidRequest(err, "emails")
		}
		return zList, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	// get list friends from userId
	friends, err := h.repo.GetFriendshipByUserIDAndStatus(ctx, userIDs[email], domain.FriendshipStatusFriended)
	if err != nil {
		return zList, common.ErrCannotListEntity(domain.Friendship{}.DomainName(), err)
	}

	// get friendIDs from userId or friendId field if not same userID
	r := make([]string, 0, len(friends))
	for _, v := range friends {
		if v.UserID == userIDs[email] {
			r = append(r, v.FriendID)
		} else {
			r = append(r, v.UserID)
		}
	}

	// get email from userID
	friendIDs, err := h.userRepo.GetEmailsByUserIDs(ctx, r)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		return zList, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	result := make([]string, 0, len(friendIDs))
	for _, v:= range friendIDs{
		result = append(result, v)
	}

	return result, nil
}
