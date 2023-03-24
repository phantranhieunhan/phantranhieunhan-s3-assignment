package query

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
)

type ListUpdatesUserHandler struct {
	userRepo         domain.UserRepo
	subscriptionRepo domain.SubscriptionRepo
}

func NewListUpdatesUserHandler(subscriptionRepo domain.SubscriptionRepo, userRepo domain.UserRepo) ListUpdatesUserHandler {
	return ListUpdatesUserHandler{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
	}
}

func (h ListUpdatesUserHandler) Handle(ctx context.Context, email, text string) ([]string, error) {
	emailFromTexts := util.GetEmailsFromString(text)
	emails := append(emailFromTexts, email)
	emails = util.RemoveDuplicates(emails)

	// get userId from email to check available
	mapEmailUser, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return nil, common.ErrInvalidRequest(err, "emails")
		}
		return []string{}, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}
	userID, ok := mapEmailUser[email]
	if !ok {
		return []string{}, common.ErrInvalidRequest(nil, "email")
	}
	// get list subscription from userId
	subs, err := h.subscriptionRepo.GetSubscriptionEmailsByUserIDAndStatus(ctx, userID, domain.SubscriptionStatusSubscribed)
	if err != nil {
		logger.Errorf("subscriptionRepo.GetSubscriptionEmailsByUserIDAndStatus %w", err)
		return []string{}, common.ErrCannotListEntity(domain.Subscription{}.DomainName(), err)
	}

	subs = append(subs, emailFromTexts...)
	subs = util.RemoveDuplicates(subs)

	result := subs
	if len(emailFromTexts) > 0 {
		blockSubs, err := h.subscriptionRepo.GetSubscriptionEmailsByUserIDAndStatus(ctx, userID, domain.SubscriptionStatusUnsubscribed)
		if err != nil {
			logger.Errorf("subscriptionRepo.GetSubscriptionEmailsByUserIDAndStatus %w", err)
			return []string{}, common.ErrCannotListEntity(domain.Subscription{}.DomainName(), err)
		}

		if len(blockSubs) > 0 {
			result = make([]string, 0)
			for _, r := range subs {
				if !util.IsContain(blockSubs, r) {
					result = append(result, r)
				}
			}
		}
	}

	return result, nil
}
