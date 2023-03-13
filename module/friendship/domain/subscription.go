package domain

type SubscriptionStatus int

const (
	SubscriptionStatusInvalid SubscriptionStatus = iota
	SubscriptionStatusSubscribed
	SubscriptionStatusUnsubscribed
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

type Subscriptions []Subscription
