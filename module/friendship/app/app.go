package app

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command/payload"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ConnectFriendship interface {
		Handle(ctx context.Context, userEmail string, friendEmail string) (domain.Friendship, error)
	}
	SubscribeUser interface {
		Handle(ctx context.Context, payload payload.SubscriberUserPayloads) error
		HandleWithSubscription(ctx context.Context, ds domain.Subscriptions) error
	}
	BlockUpdatesUser interface {
		Handle(ctx context.Context, payload payload.BlockUpdatesUserPayload) error
	}
}

type Queries struct {
	ListFriends interface {
		Handle(ctx context.Context, email string) ([]string, error)
	}
	ListCommonFriends interface {
		Handle(ctx context.Context, emails []string) ([]string, error)
	}
	ListUpdatesUser interface {
		Handle(ctx context.Context, email string, text string) ([]string, error)
	}
	ListSubscriptions interface {
		Handle(ctx context.Context) ([]domain.FullSubscription, error)
	}
}
