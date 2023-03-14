package command

import (
	"context"
	"errors"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type SubscribeUserRepo interface {
	Create(ctx context.Context, sub domain.Subscription) (string, error)
	GetSubscription(ctx context.Context, ss domain.Subscriptions) (domain.Subscriptions, error)
	UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error
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

func (h SubscribeUserHandler) Handle(ctx context.Context, payload SubscriberUserPayloads) error {
	emails := payload.GetEmails()
	if len(emails) < 2 {
		return common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "payload")
	}
	userIDs, _, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return common.ErrInvalidRequest(err, "emails")
		}
		return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
	}

	ds := make(domain.Subscriptions, 0, len(payload))
	mapSub := make(map[string]domain.Subscription, len(payload))

	for _, v := range payload {
		targetID := userIDs[v.Target]
		requestorID := userIDs[v.Requestor]
		sc := domain.Subscription{
			UserID:       targetID,
			SubscriberID: requestorID,
		}
		ds = append(ds, sc)
		mapSub[sc.GetMapKey()] = sc
	}
	return h.handle(ctx, ds, mapSub)
}

func (h SubscribeUserHandler) HandleWithSubscription(ctx context.Context, ds domain.Subscriptions) error {
	mapSub := make(map[string]domain.Subscription, len(ds))
	for _, v := range ds {
		mapSub[v.GetMapKey()] = v
	}

	return h.handle(ctx, ds, mapSub)
}

func (h SubscribeUserHandler) handle(ctx context.Context, ds domain.Subscriptions, mapSub map[string]domain.Subscription) error {
	err := h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		fs, err := h.subscribeUserRepo.GetSubscription(ctx, ds)
		if err != nil {
			logger.Errorf("subscribeUserRepo.GetSubscription %w", err)
			return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
		}

		for _, v := range fs {
			mapSub[v.GetMapKey()] = v
		}
		for _, v := range mapSub {
			if v.Status == domain.SubscriptionStatusUnsubscribed || v.Status == domain.SubscriptionStatusInvalid {
				f, err := h.friendshipRepo.GetFriendshipByUserIDs(ctx, v.UserID, v.SubscriberID)
				if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
					return common.ErrCannotGetEntity(f.DomainName(), err)
				}
				if f.Status.CanNotSubscribe() {
					return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
				}
				if v.Status == domain.SubscriptionStatusInvalid {
					v.Status = domain.SubscriptionStatusSubscribed
					v.Id, err = h.subscribeUserRepo.Create(ctx, v)
					if err != nil {
						return common.ErrCannotCreateEntity(v.DomainName(), err)
					}
				} else {
					if err := h.subscribeUserRepo.UpdateStatus(ctx, v.Id, domain.SubscriptionStatusSubscribed); err != nil {
						return common.ErrCannotUpdateEntity(v.DomainName(), err)
					}
				}
			}
		}
		return nil
	})

	return err
}
