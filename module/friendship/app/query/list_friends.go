package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListFriendsHandler struct {
	repo     domain.FriendshipRepo
	userRepo domain.UserRepo
}

func NewListFriendsHandler(repo domain.FriendshipRepo, userRepo domain.UserRepo) ListFriendsHandler {
	return ListFriendsHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h ListFriendsHandler) Handle(ctx context.Context, email string) ([]string, error) {
	// get userId from email to check available
	mapEmailUser, err := h.userRepo.GetUserIDsByEmails(ctx, []string{email})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return nil, common.ErrInvalidRequest(err, "emails")
		}
		return nil, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	// get list friends from userId
	result, err := h.repo.GetFriendshipByUserIDAndStatus(ctx, mapEmailUser, domain.FriendshipStatusFriended)
	if err != nil {
		logger.Errorf("friendshipRepo.GetFriendshipByUserIDAndStatus %w", err)
		return nil, common.ErrCannotListEntity(domain.Friendship{}.DomainName(), err)
	}

	return result, nil
}
