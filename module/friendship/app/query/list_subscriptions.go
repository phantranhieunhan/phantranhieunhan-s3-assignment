package query

import (
	"context"
	"fmt"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListSubscriptionsHandler struct {
	subscriptionRepo domain.SubscriptionRepo
}

func NewListSubscriptionsHandler(repo domain.SubscriptionRepo) ListSubscriptionsHandler {
	return ListSubscriptionsHandler{
		subscriptionRepo: repo,
	}
}

func (l ListSubscriptionsHandler) Handle(ctx context.Context) ([]domain.FullSubscription, error) {
	list, err := l.subscriptionRepo.GetAll(ctx)
	if err != nil {
		return []domain.FullSubscription{}, fmt.Errorf("Error subscriptionRepo.GetAll %w", err)
	}
	return list, nil
}
