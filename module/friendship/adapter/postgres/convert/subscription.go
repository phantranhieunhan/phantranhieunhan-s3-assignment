package convert

import (
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

func ToSubscriptionModel(d domain.Subscription) model.Subscription {
	return model.Subscription{
		ID:           d.Id,
		UserID:       d.UserID,
		SubscriberID: d.SubscriberID,
		Status:       int(d.Status),
	}
}

func ToSubscriptionDomain(d model.Subscription) domain.Subscription {
	return domain.Subscription{
		Base: domain.Base{
			Id:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		UserID:       d.UserID,
		SubscriberID: d.SubscriberID,
		Status:       domain.SubscriptionStatus(d.Status),
	}
}

func ToSubscriptionsDomain(m model.SubscriptionSlice) domain.Subscriptions {
	ds := make(domain.Subscriptions, 0, len(m))
	for _, v := range m {
		ds = append(ds, ToSubscriptionDomain(*v))
	}
	return ds
}
