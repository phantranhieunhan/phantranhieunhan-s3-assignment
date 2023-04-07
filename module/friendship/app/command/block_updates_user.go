package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command/payload"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

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

func (b BlockUpdatesUserHandler) Handle(ctx context.Context, payload payload.BlockUpdatesUserPayload) error {
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
		gotSub, err := b.subscriptionRepo.GetSubscription(ctx, domain.Subscriptions{
			{UserID: targetID, SubscriberID: requestorID},
		})
		if err != nil {
			logger.Errorf("subscribeUserRepo.GetSubscription %w", err)
			return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
		}
		if len(gotSub) > 0 {
			if !gotSub[0].Status.AllowBlock() {
				return common.ErrInvalidRequest(domain.ErrAlreadyExists, "emails")
			}
		}

		f, err := b.friendshipRepo.GetFriendshipByUserIDs(ctx, requestorID, targetID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(f.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound || f.Status.CanBlockUser() {
			if err = b.blockUser(ctx, f.Id, requestorID, targetID); err != nil {
				return err
			}
		}

		if err = b.unsubscribeUser(ctx, requestorID, targetID); err != nil {
			return err
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
	} else if err := b.friendshipRepo.UpdateStatus(ctx, friendshipID, domain.FriendshipStatusBlocked); err != nil {
		logger.Errorf("repo.UpdateStatus %w", err)
		return common.ErrCannotUpdateEntity(domain.Friendship{}.DomainName(), err)
	}

	return nil
}

func (b BlockUpdatesUserHandler) unsubscribeUser(ctx context.Context, requestorID, targetID string) error {
	sub := domain.Subscription{UserID: targetID, SubscriberID: requestorID, Status: domain.SubscriptionStatusUnsubscribed}
	_, err := b.subscriptionRepo.UpsertSubscription(ctx, sub)
	if err != nil {
		logger.Errorf("subscriptionRepo.UpsertSubscription %w", err)
		return common.ErrCannotUpdateEntity(sub.DomainName(), err)
	}

	return nil
}
