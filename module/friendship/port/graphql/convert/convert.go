package convert

import (
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/graphql/model"
)

func ToSubscriptionG(list []domain.FullSubscription) (result []*model.SubscriptionG) {
	for _, v := range list {
		result = append(result, &model.SubscriptionG{
			ID:         v.Id,
			CreatedAt:  float64(v.CreatedAt.UnixMilli()),
			UpdatedAt:  float64(v.UpdatedAt.UnixMilli()),
			User:       &model.User{ID: v.UserID, Email: v.User.Email},
			Subscriber: &model.User{ID: v.SubscriberID, Email: v.Subscriber.Email},
			Status:     model.AllSubscriptionStatus[v.Status],
		})
	}
	return
}
