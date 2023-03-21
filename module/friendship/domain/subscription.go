package domain

import "errors"

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

func (s SubscriptionStatus) IsNoneExisted() bool {
	switch s {
	case SubscriptionStatusInvalid:
		return true
	default:
		return false
	}
}

var (
	ErrCannotCreateSubscription = errors.New("cannot create subscription")
	ErrNeedAtLeastTwoEmails     = errors.New("need at least two emails")
)

type Subscription struct {
	Base         `json:",inline"`
	UserID       string             `json:"user_id"`
	SubscriberID string             `json:"subscriber_id"`
	Status       SubscriptionStatus `json:"status"`
}

func (r Subscription) DomainName() string {
	return "Subscription"
}

func (r Subscription) GetMapKey() string {
	return r.UserID + r.SubscriberID
}

type Subscriptions []Subscription
