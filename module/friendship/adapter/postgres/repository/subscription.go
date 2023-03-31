package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/view"
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
	if err := m.Insert(ctx, s.db.Model(ctx), boil.Infer()); err != nil {
		return "", common.ErrDB(err)
	}
	return m.ID, nil
}

func (f SubscriptionRepository) UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error {
	m := model.Subscription{
		ID:     id,
		Status: int(status),
	}
	effectedRows, err := m.Update(ctx, f.db.Model(ctx), boil.Whitelist(model.SubscriptionColumns.Status, model.SubscriptionColumns.UpdatedAt))
	if err != nil {
		return common.ErrDB(err)
	}

	if effectedRows != 1 {
		return common.ErrDB(domain.ErrUpdateRecordNotFound)
	}

	return nil
}

func (f SubscriptionRepository) UnsertSubscription(ctx context.Context, sub domain.Subscription) (string, error) {
	m := convert.ToSubscriptionModel(sub)
	if m.ID == "" {
		m.ID = util.GenUUID()
	}
	conflictFields := []string{model.SubscriptionColumns.UserID, model.SubscriptionColumns.SubscriberID}
	err := m.Upsert(ctx, f.db.Model(ctx), true, conflictFields, boil.Whitelist(model.SubscriptionColumns.Status, model.FriendshipColumns.UpdatedAt), boil.Infer())
	if err != nil {
		return "", common.ErrDB(err)
	}
	return m.ID, nil
}

func (s SubscriptionRepository) GetSubscription(ctx context.Context, ss domain.Subscriptions) (domain.Subscriptions, error) {
	where := make([]qm.QueryMod, 0)
	for _, v := range ss {
		where = append(where, qm.Or("user_id = ? AND subscriber_id = ?", v.UserID, v.SubscriberID))
	}

	m, err := model.Subscriptions(where...).All(ctx, s.db.Model(ctx))
	if err != nil {
		return domain.Subscriptions{}, common.ErrDB(err)
	}
	return convert.ToSubscriptionsDomain(m), nil
}

func (s SubscriptionRepository) GetSubscriptionEmailsByUserIDAndEmails(ctx context.Context, id string, emails []string) ([]string, error) {
	list := make([]view.SubscriberEmail, 0)
	query := `select distinct * from (
		select email from public.users u1 where id in (select subscriber_id from subscriptions where user_id = $1 and status = $2)
		union
		select email from public.users u2
		where email = any('{%s}'::text[])
		and id not in (select subscriber_id from subscriptions where user_id = $1 and status = $3)
	) as sub_query`
	iEmails := strings.Join(emails, ",")

	queryWithEmails := fmt.Sprintf(query, iEmails)
	err := model.NewQuery(
		qm.SQL(queryWithEmails, id, domain.SubscriptionStatusSubscribed, domain.SubscriptionStatusUnsubscribed),
	).Bind(ctx, s.db.Model(ctx), &list)
	if err != nil {
		return []string{}, common.ErrDB(err)
	}

	result := make([]string, len(list))
	for i := range result {
		result[i] = list[i].Email
	}

	return result, nil
}
