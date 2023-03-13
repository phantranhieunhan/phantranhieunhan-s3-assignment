package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListCommonFriendsRepo interface {
	GetFriendshipByUserIDAndStatus(ctx context.Context, userID string, status ...domain.FriendshipStatus) (domain.Friendships, error)
}

type ListCommonFriendsHandler struct {
	repo     ListCommonFriendsRepo
	userRepo UserRepo
}

func NewListCommonFriendsHandler(repo ListCommonFriendsRepo, userRepo UserRepo) ListCommonFriendsHandler {
	return ListCommonFriendsHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (h ListCommonFriendsHandler) Handle(ctx context.Context, emails []string) ([]string, error) {
	emptyList := []string{}
	if len(emails) != 2 {
		return emptyList, common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "emails")
	}
	// get userId from email to check available
	mapUserIDs, userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return emptyList, common.ErrInvalidRequest(err, "emails")
		}
		return emptyList, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}
	f := make([]string, 0)

	for _, e := range emails {
		// get list friends from userId
		friends, err := h.repo.GetFriendshipByUserIDAndStatus(ctx, mapUserIDs[e], domain.FriendshipStatusFriended)
		if err != nil {
			logger.Errorf("friendshipRepo.GetFriendshipByUserIDAndStatus %w", err)
			return emptyList, common.ErrCannotListEntity(domain.Friendship{}.DomainName(), err)
		}

		// get friendIDs from userId or friendId field if not same userID
		for _, v := range friends {
			var y string
			if v.UserID == mapUserIDs[e] {
				y = v.FriendID
			} else {
				y = v.UserID
			}

			// remove owner from friends list
			if !checkBlacklist(userIDs, y) {
				f = append(f, y)
			}
		}
	}

	mutual := getMutual(f)

	// get email from userID
	_, fEmails, err := h.userRepo.GetEmailsByUserIDs(ctx, mutual)
	if err != nil {
		logger.Errorf("userRepo.GetEmailsByUserIDs %w", err)
		return emptyList, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	return fEmails, nil
}

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

func checkBlacklist(blacklist []string, s string) bool {
	for _, item := range blacklist {
		if item == s {
			return true
		}
	}
	return false
}
