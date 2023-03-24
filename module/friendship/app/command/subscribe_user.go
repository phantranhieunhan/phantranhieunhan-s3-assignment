package command

import (
	"context"
	"errors"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
)

const EMAIL_TOTAL = 2

type SubscribeUserHandler struct {
	friendshipRepo    domain.FriendshipRepo
	userRepo          domain.UserRepo
	subscribeUserRepo domain.SubscriptionRepo
	transactor        Transactor
}

func NewSubscribeUserHandler(repo domain.FriendshipRepo, userRepo domain.UserRepo, subscribeUserRepo domain.SubscriptionRepo, transactor Transactor) SubscribeUserHandler {
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
	userIds := make([]string, 0, len(s)*EMAIL_TOTAL)
	for _, u := range s {
		if !util.IsContain(userIds, u.Requestor) {
			userIds = append(userIds, u.Requestor)
		}
		if !util.IsContain(userIds, u.Target) {
			userIds = append(userIds, u.Target)
		}
	}

	return userIds
}

func (h SubscribeUserHandler) Handle(ctx context.Context, payload SubscriberUserPayloads) error {
	emails := payload.GetEmails()
	if len(emails) < EMAIL_TOTAL {
		return common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "payload")
	}
	userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, emails)
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return common.ErrInvalidRequest(err, "emails")
		}
		return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
	}

	ds := make(domain.Subscriptions, 0, len(payload))

	for _, v := range payload {
		sc := domain.Subscription{
			UserID:       userIDs[v.Target],
			SubscriberID: userIDs[v.Requestor],
		}
		ds = append(ds, sc)
	}
	return h.handle(ctx, ds)
}

func (h SubscribeUserHandler) HandleWithSubscription(ctx context.Context, ds domain.Subscriptions) error {
	return h.handle(ctx, ds)
}

func (h SubscribeUserHandler) handle(ctx context.Context, ds domain.Subscriptions) error {
	mapSub := make(map[string]domain.Subscription, 0)
	for _, v := range ds {
		mapSub[v.GetMapKey()] = v
	}

	err := h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		fs, err := h.subscribeUserRepo.GetSubscription(ctx, ds)
		if err != nil {
			logger.Errorf("subscribeUserRepo.GetSubscription %w", err)
			return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
		}

		for _, v := range fs {
			mapSub[v.GetMapKey()] = v
		}
		for _, sub := range mapSub {
			if sub.Status.AllowSubscribe() {
				friendship, err := h.friendshipRepo.GetFriendshipByUserIDs(ctx, sub.UserID, sub.SubscriberID)
				if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
					return common.ErrCannotGetEntity(friendship.DomainName(), err)
				}

				if friendship.Status.CanNotSubscribe() {
					return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
				}

				if sub.Status.IsNoneExisted() {
					sub.Status = domain.SubscriptionStatusSubscribed
					sub.Id, err = h.subscribeUserRepo.Create(ctx, sub)
					if err != nil {
						return common.ErrCannotCreateEntity(sub.DomainName(), err)
					}
				} else {
					if err := h.subscribeUserRepo.UpdateStatus(ctx, sub.Id, domain.SubscriptionStatusSubscribed); err != nil {
						return common.ErrCannotUpdateEntity(sub.DomainName(), err)
					}
				}
			}
		}
		return nil
	})

	return err
}
