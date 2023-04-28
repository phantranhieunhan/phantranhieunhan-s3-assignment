package view

import (
	"time"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type SubscriberEmail struct {
	Email string `boil:"email"`
}

type FullSubscription struct {
	Id        string    `json:"id" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID          string                    `json:"user_id"`
	SubscriberID    string                    `json:"subscriber_id"`
	Status          domain.SubscriptionStatus `json:"status"`
	UserEmail       string                    `json:"user_email,omitempty"`
	SubscriberEmail string                    `json:"subscriber_email,omitempty"`
}
