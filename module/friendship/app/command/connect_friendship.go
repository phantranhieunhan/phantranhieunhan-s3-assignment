package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ConnectFriendshipHandler struct {
	friendshipRepo       domain.FriendshipRepo
	userRepo             domain.UserRepo
	transactor           Transactor
	subscribeUserCommand domain.SubscribeUserCommand
}

func NewConnectFriendshipHandler(repo domain.FriendshipRepo, userRepo domain.UserRepo, transactor Transactor, subscribeUser domain.SubscribeUserCommand) ConnectFriendshipHandler {
	return ConnectFriendshipHandler{
		friendshipRepo:       repo,
		userRepo:             userRepo,
		transactor:           transactor,
		subscribeUserCommand: subscribeUser,
	}
}

func (h ConnectFriendshipHandler) Handle(ctx context.Context, userEmail, friendEmail string) (string, error) {
	userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, []string{userEmail, friendEmail})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return "", common.ErrInvalidRequest(err, "emails")
		}
		return "", common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}
	var id string
	d := domain.Friendship{
		Status:   domain.FriendshipStatusFriended,
		UserID:   userIDs[userEmail],
		FriendID: userIDs[friendEmail],
	}

	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := h.friendshipRepo.GetFriendshipByUserIDs(ctx, d.UserID, d.FriendID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(d.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound {
			id, err = h.friendshipRepo.Create(ctx, d)
			if err != nil {
				logger.Errorf("repo.Create %w", err)
				return common.ErrCannotCreateEntity(d.DomainName(), err)
			}
		} else {
			if !f.Status.CanConnect() {
				logger.Errorf("Status.CanConnect")
				return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
			}
			if err = h.friendshipRepo.UpdateStatus(ctx, f.Id, domain.FriendshipStatusFriended); err != nil {
				logger.Errorf("repo.UpdateStatus %w", err)
				return common.ErrCannotUpdateEntity(d.DomainName(), err)
			}
			id = f.Id
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	err = h.subscribeUserCommand.HandleWithSubscription(ctx, domain.Subscriptions{
		domain.Subscription{
			UserID:       d.UserID,
			SubscriberID: d.FriendID,
		},
		domain.Subscription{
			UserID:       d.FriendID,
			SubscriberID: d.UserID,
		},
	})
	if err != nil {
		logger.Errorf("Create Subscription fail when create connection friendship HandleWithSubscription: %w", err)
		return "", err
	}

	return id, err
}
