package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type FriendshipRepo interface {
	Create(ctx context.Context, d domain.Friendship) (string, error)
	GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error
}

type UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, []string, error)
}

type SubscribeUserMQ interface {
	SubscribeUser(ctx context.Context, ds domain.Subscriptions) error
}

type ConnectFriendshipHandler struct {
	friendshipRepo   FriendshipRepo
	userRepo         UserRepo
	subscriptionRepo SubscribeUserRepo
	transactor       Transactor
	subscribeUserMQ  SubscribeUserMQ
}

func NewConnectFriendshipHandler(repo FriendshipRepo, userRepo UserRepo, subRepo SubscribeUserRepo, transactor Transactor, subMq SubscribeUserMQ) ConnectFriendshipHandler {
	return ConnectFriendshipHandler{
		friendshipRepo:   repo,
		userRepo:         userRepo,
		subscriptionRepo: subRepo,
		transactor:       transactor,
		subscribeUserMQ:  subMq,
	}
}

func (h ConnectFriendshipHandler) Handle(ctx context.Context, userEmail, friendEmail string) (string, error) {
	userIDs, _, err := h.userRepo.GetUserIDsByEmails(ctx, []string{userEmail, friendEmail})
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

	err = h.subscribeUserMQ.SubscribeUser(ctx, domain.Subscriptions{
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
