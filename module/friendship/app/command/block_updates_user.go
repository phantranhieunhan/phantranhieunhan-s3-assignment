package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type BlockUpdatesUserPayload struct {
	Requestor string
	Target    string
}

type BlockUpdatesUserHandler struct {
	friendshipRepo   domain.FriendshipRepo
	userRepo         domain.UserRepo
	subscriptionRepo domain.SubscriptionRepo
	transactor       Transactor
}

func NewBlockUpdatesUserHandler(repo domain.FriendshipRepo, userRepo domain.UserRepo, subRepo domain.SubscriptionRepo, transactor Transactor) BlockUpdatesUserHandler {
	return BlockUpdatesUserHandler{
		friendshipRepo:   repo,
		userRepo:         userRepo,
		subscriptionRepo: subRepo,
		transactor:       transactor,
	}
}

func (b BlockUpdatesUserHandler) Handle(ctx context.Context, payload BlockUpdatesUserPayload) error {
	if payload.Requestor == payload.Target {
		return common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "payload")
	}

	userIDs, err := b.userRepo.GetUserIDsByEmails(ctx, []string{payload.Requestor, payload.Target})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return common.ErrInvalidRequest(err, "emails")
		}
		return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
	}

	requestorID := userIDs[payload.Requestor]
	targetID := userIDs[payload.Target]

	err = b.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := b.friendshipRepo.GetFriendshipByUserIDs(ctx, requestorID, targetID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(f.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound || f.Status.CanBlockUser() {
			if err = b.unsubscribeUser(ctx, requestorID, targetID); err != nil {
				return err
			}

			if err = b.blockUser(ctx, f.Id, requestorID, targetID); err != nil {
				return err
			}
		} else if f.Status == domain.FriendshipStatusFriended {
			if err = b.unsubscribeUser(ctx, requestorID, targetID); err != nil {
				return err
			}
		} else {
			return common.ErrInvalidRequest(domain.ErrCannotBlockUpdatesFromBlockedUser, "")
		}

		return nil
	})
	return err
}

func (b BlockUpdatesUserHandler) blockUser(ctx context.Context, friendshipID, requestorID, targetID string) error {
	if friendshipID == "" {
		d := domain.Friendship{}.FriendshipWithBlock(requestorID, targetID)
		_, err := b.friendshipRepo.Create(ctx, d)
		if err != nil {
			logger.Errorf("repo.Create %w", err)
			return common.ErrCannotCreateEntity(d.DomainName(), err)
		}
	} else {
		if err := b.friendshipRepo.UpdateStatus(ctx, friendshipID, domain.FriendshipStatusBlocked); err != nil {
			logger.Errorf("repo.UpdateStatus %w", err)
			return common.ErrCannotUpdateEntity(domain.Friendship{}.DomainName(), err)
		}
	}

	return nil
}

func (b BlockUpdatesUserHandler) unsubscribeUser(ctx context.Context, requestorID, targetID string) error {
	sub := domain.Subscription{UserID: targetID, SubscriberID: requestorID, Status: domain.SubscriptionStatusUnsubscribed}
	err := b.subscriptionRepo.UnsertSubscription(ctx, sub)
	if err != nil {
		return common.ErrCannotUpdateEntity(sub.DomainName(), err)
	}

	return nil
}
