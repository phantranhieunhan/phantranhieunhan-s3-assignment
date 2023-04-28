package domain

import (
	"context"
	"errors"
)

type SubscriptionStatus int

const (
	SubscriptionStatusInvalid SubscriptionStatus = iota
	SubscriptionStatusSubscribed
	SubscriptionStatusUnsubscribed
)

func (s SubscriptionStatus) AllowSubscribe() bool {
	switch s {
	case SubscriptionStatusInvalid, SubscriptionStatusUnsubscribed:
		return true
	default:
		return false
	}
}

func (s SubscriptionStatus) AllowBlock() bool {
	switch s {
	case SubscriptionStatusInvalid, SubscriptionStatusSubscribed:
		return true
	default:
		return false
	}
}

func (s SubscriptionStatus) IsNoneExisted() bool {
	switch s {
	case SubscriptionStatusInvalid:
		return true
	default:
		return false
	}
}

var (
	ErrCannotCreateSubscription          = errors.New("cannot create subscription")
	ErrNeedAtLeastTwoEmails              = errors.New("need at least two emails")
	ErrCannotBlockUpdatesFromBlockedUser = errors.New("cannot block updates from blocked user")
)

type Subscription struct {
	Base         `json:",inline"`
	UserID       string             `json:"user_id"`
	SubscriberID string             `json:"subscriber_id"`
	Status       SubscriptionStatus `json:"status"`
}

type FullSubscription struct {
	Subscription
	User       User
	Subscriber User
}

func (r Subscription) DomainName() string {
	return "Subscription"
}

func (r Subscription) GetUserSubscriberMapKey() string {
	return r.UserID + r.SubscriberID
}

func (r Subscription) GetSubscriberUserMapKey() string {
	return r.SubscriberID + r.UserID
}

type Subscriptions []Subscription

type SubscribeUserCommand interface {
	HandleWithSubscription(ctx context.Context, ds Subscriptions) error
}

type SubscriptionRepo interface {
	Create(ctx context.Context, sub Subscription) (string, error)
	GetSubscription(ctx context.Context, ss Subscriptions) (Subscriptions, error)
	UpdateStatus(ctx context.Context, id string, status SubscriptionStatus) error
	UpsertSubscription(ctx context.Context, sub Subscription) (string, error)
	GetSubscriptionEmailsByUserIDAndEmails(ctx context.Context, id string, emails []string) ([]string, error)
	GetAll(ctx context.Context) ([]FullSubscription, error)
}
