package repository

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SubscriptionRepository struct {
	db postgres.Database
}

func NewSubscriptionRepository(db postgres.Database) SubscriptionRepository {
	return SubscriptionRepository{
		db: db,
	}
}

func (s SubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (string, error) {
	sub.Id = util.GenUUID()
	m := convert.ToSubscriptionModel(sub)
	if err := m.Insert(ctx, s.db.DB, boil.Infer()); err != nil {
		return "", common.ErrDB(err)
	}
	return m.ID, nil
}

func (f SubscriptionRepository) UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error {
	m := model.Subscription{
		ID:     id,
		Status: int(status),
	}
	_, err := m.Update(ctx, f.db.DB, boil.Whitelist(model.SubscriptionColumns.Status, model.SubscriptionColumns.UpdatedAt))
	if err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s SubscriptionRepository) GetSubscription(ctx context.Context, ss domain.Subscriptions) (domain.Subscriptions, error) {
	where := make([]qm.QueryMod, 0)
	for _, v := range ss {
		where = append(where, qm.Or("user_id = ? AND subscriber_id = ?", v.UserID, v.SubscriberID))
	}

	m, err := model.Subscriptions(where...).All(ctx, s.db.DB)
	if err != nil {
		return domain.Subscriptions{}, common.ErrDB(err)
	}
	return convert.ToSubscriptionsDomain(m), nil
}
