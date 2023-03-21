package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

const EMAIL_TOTAL = 2

type ListCommonFriends_FriendshipRepo interface {
	GetFriendshipByUserIDAndStatus(ctx context.Context, mapEmailUser map[string]string, status ...domain.FriendshipStatus) ([]string, error)
}

type ListCommonFriends_UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}

type ListCommonFriendsHandler struct {
	repo     ListCommonFriends_FriendshipRepo
	userRepo ListCommonFriends_UserRepo
}

func NewListCommonFriendsHandler(repo ListCommonFriends_FriendshipRepo, userRepo ListCommonFriends_UserRepo) ListCommonFriendsHandler {
	return ListCommonFriendsHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h ListCommonFriendsHandler) Handle(ctx context.Context, emails []string) ([]string, error) {
	if len(emails) != EMAIL_TOTAL {
		return nil, common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "emails")
	}
	// get userId from email to check available
	mapEmailUserIDs, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return nil, common.ErrInvalidRequest(err, "emails")
		}
		return nil, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	friends, err := h.repo.GetFriendshipByUserIDAndStatus(ctx, mapEmailUserIDs, domain.FriendshipStatusFriended)
	if err != nil {
		logger.Errorf("friendshipRepo.GetFriendshipByUserIDAndStatus %w", err)
		return nil, common.ErrCannotListEntity(domain.Friendship{}.DomainName(), err)
	}

	mutual := getMutual(friends)

	return mutual, nil
}

// after add all friends of 2 user in to a list, then get items is duplicated
func getMutual(fullList []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range fullList {
		if _, value := allKeys[item]; value {
			list = append(list, item)
		} else {
			allKeys[item] = true
		}
	}
	return list
}
