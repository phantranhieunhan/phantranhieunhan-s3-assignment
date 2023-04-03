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
	emailFromTexts := util.RemoveDuplicates(util.GetEmailsFromString(text))

	// get userId from email to check available
	mapEmailUser, err := h.userRepo.GetUserIDsByEmails(ctx, []string{email})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return nil, common.ErrInvalidRequest(err, "emails")
		}
		return nil, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}

	userID, ok := mapEmailUser[email]
	if !ok {
		return nil, common.ErrInvalidRequest(nil, "email")
	}

	// get list subscription from userId
	subs, err := h.subscriptionRepo.GetSubscriptionEmailsByUserIDAndEmails(ctx, userID, emailFromTexts)
	if err != nil {
		logger.Errorf("subscriptionRepo.GetSubscriptionEmailsByUserIDAndStatus %w", err)
		return nil, common.ErrCannotListEntity(domain.Subscription{}.DomainName(), err)
	}

	return subs, nil
}
