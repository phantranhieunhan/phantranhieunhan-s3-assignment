package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListFriends_FriendshipRepo interface {
	GetFriendshipByUserIDAndStatus(ctx context.Context, userIDs, emails []string, status ...domain.FriendshipStatus) ([]string, error)
}

type ListFriends_UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}

type ListFriendsHandler struct {
	repo     ListFriends_FriendshipRepo
	userRepo ListFriends_UserRepo
}

func NewListFriendsHandler(repo ListFriends_FriendshipRepo, userRepo ListFriends_UserRepo) ListFriendsHandler {
	return ListFriendsHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h ListFriendsHandler) Handle(ctx context.Context, email string) ([]string, error) {
	var emptyList []string

	// get userId from email to check available
	userIDs, _, err := h.userRepo.GetUserIDsByEmails(ctx, []string{email})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return emptyList, common.ErrInvalidRequest(err, "emails")
		}
		return emptyList, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	// get list friends from userId
	result, err := h.repo.GetFriendshipByUserIDAndStatus(ctx, []string{userIDs[email]}, []string{email}, domain.FriendshipStatusFriended)
	if err != nil {
		logger.Errorf("friendshipRepo.GetFriendshipByUserIDAndStatus %w", err)
		return emptyList, common.ErrCannotListEntity(domain.Friendship{}.DomainName(), err)
	}

	return result, nil
}
