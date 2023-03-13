package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type SubscribeUserRepo interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	GetSubscription(ctx context.Context, ss domain.Subscriptions) (domain.Subscriptions, error)
}

type SubscribeUserHandler struct {
	friendshipRepo    ConnectFriendshipRepo
	userRepo          UserRepo
	subscribeUserRepo SubscribeUserRepo
	transactor        Transactor
}

func NewSubscribeUserHandler(repo ConnectFriendshipRepo, userRepo UserRepo, subscribeUserRepo SubscribeUserRepo, transactor Transactor) SubscribeUserHandler {
	return SubscribeUserHandler{
		friendshipRepo:    repo,
		userRepo:          userRepo,
		subscribeUserRepo: subscribeUserRepo,
		transactor:        transactor,
	}
}

type SubscriberUserPayload struct {
	Requestor string
	Target    string
}

type SubscriberUserPayloads []SubscriberUserPayload

func (s SubscriberUserPayloads) GetEmails() []string {
	userIds := make([]string, 0, len(s)*2)
	for _, u := range s {
		if !CheckExisted(userIds, u.Requestor) {
			userIds = append(userIds, u.Requestor)
		}
		if !CheckExisted(userIds, u.Target) {
			userIds = append(userIds, u.Target)
		}
	}

	return userIds
}

func CheckExisted(list []string, p string) bool {
	for _, v := range list {
		if v == p {
			return true
		}
	}
	return false
}

func (h SubscribeUserHandler) Handle(ctx context.Context, payload SubscriberUserPayloads) (string, error) {
	emails := payload.GetEmails()
	if len(emails) < 2 {
		// return error
	}
	userIDs, _, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return "", common.ErrInvalidRequest(err, "emails")
		}
		return "", common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
	}

	ds := make(domain.Subscriptions, 0, len(payload))

	for _, v := range payload {
		ds = append(ds, domain.Subscription{
			UserID:       v.Target,
			SubscriberID: v.Requestor,
		})
	}

	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := h.subscribeUserRepo.GetSubscription(ctx, ds)
		if err != nil {
			logger.Errorf("subscribeUserRepo.GetSubscription %w", err)
			return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
		}
		createdSubscription := make(sub)
	})

	return "", nil
}
