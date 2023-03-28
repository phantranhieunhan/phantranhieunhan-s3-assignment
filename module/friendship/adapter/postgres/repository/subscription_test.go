package repository

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Suite struct {
	db postgres.Database
}

func NewSuite(ctx context.Context) Suite {
	err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
	}
	db := postgres.NewDatabase()
	return Suite{db: db}
}

type SubscriptionTestCase struct {
	name      string
	modifySub domain.Subscription
	isExisted bool
	isFounded bool
	err       error
}

func TestSubscription_Create(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	errDB := errors.New("some error from db")

	cs := []SubscriptionTestCase{
		{
			name: "successful",
		},
		{
			name: "fail by reference user id",
			modifySub: domain.Subscription{
				UserID: "user-id",
			},
			err: common.ErrDB(errDB),
		},
		{
			name: "fail by reference subscriber id",
			modifySub: domain.Subscription{
				SubscriberID: "subscriber-id",
			},
			err: common.ErrDB(errDB),
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusSubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			createdSub := sub
			if tc.err != nil {
				if tc.modifySub.SubscriberID != "" {
					createdSub.SubscriberID = "subscriber-id"
				} else if tc.modifySub.UserID != "" {
					createdSub.UserID = "user-id"
				}
			}
			var err error
			sub.Id, err = repo.Create(ctx, createdSub)
			if tc.err == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, sub.Id)
			} else {
				assert.Error(t, err)
				assert.Empty(t, sub.Id)
			}

			suite.rollbackSubscription(t, ctx, sub)
		})

	}
}

func TestSubscription_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	errDB := errors.New("some error from db")

	cs := []SubscriptionTestCase{
		{
			name: "successful",
		},
		{
			name: "fail by invalid id",
			modifySub: domain.Subscription{
				Base: domain.Base{Id: "sub-id"},
			},
			err: common.ErrDB(errDB),
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusSubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			var err error
			sub.Id, err = repo.Create(ctx, sub)
			assert.NoError(t, err)
			id := sub.Id
			if tc.modifySub.Id != "" {
				id = tc.modifySub.Id
			}

			repo.UpdateStatus(ctx, id, domain.SubscriptionStatusUnsubscribed)
			result, err := repo.GetSubscription(ctx, domain.Subscriptions{sub})
			assert.NoError(t, err)
			if tc.err == nil {
				assert.Equal(t, domain.SubscriptionStatusUnsubscribed, result[0].Status)
			} else {
				assert.Equal(t, domain.SubscriptionStatusSubscribed, result[0].Status)
			}

			suite.rollbackSubscription(t, ctx, sub)
		})
	}
}

func TestSubscription_UnsertSubscription(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	cs := []SubscriptionTestCase{
		{
			name: "successful because not existing subscription",
		},
		{
			name:      "successful because not existing subscription",
			isExisted: true,
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusSubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			var err error
			if tc.isExisted {
				sub.Id, err = repo.Create(ctx, sub)
				assert.NoError(t, err)
			}
			upsertSub := sub
			upsertSub.Status = domain.SubscriptionStatusUnsubscribed
			var upsertErr error
			sub.Id, upsertErr = repo.UnsertSubscription(ctx, upsertSub)
			assert.NoError(t, upsertErr)

			result, err := repo.GetSubscription(ctx, domain.Subscriptions{sub})
			assert.NoError(t, err)
			assert.Equal(t, upsertSub.Status, result[0].Status)

			suite.rollbackSubscription(t, ctx, sub)
		})
	}
}

func TestGetSubscription(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	cs := []SubscriptionTestCase{
		{
			name:      "successful",
			isFounded: true,
		},
		{
			name: "not found",
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusSubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			var err error
			if tc.isFounded {
				sub.Id, err = repo.Create(ctx, sub)
				assert.NoError(t, err)
			}

			result, err := repo.GetSubscription(ctx, domain.Subscriptions{
				sub,
			})
			assert.NoError(t, err)
			if tc.isFounded {
				assert.Len(t, result, 1)
			} else {
				assert.Len(t, result, 0)
			}

			suite.rollbackSubscription(t, ctx, sub)
		})

	}
}

func TestGetSubscriptionEmailsByUserIDAndStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	cs := []SubscriptionTestCase{
		{
			name:      "successful",
			isFounded: true,
		},
		{
			name: "not found",
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusUnsubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			var err error
			id := "fake-subscription"
			if tc.isFounded {
				id, err = repo.Create(ctx, sub)
				assert.NoError(t, err)
				sub.Id = id
			}

			result, err := repo.GetSubscriptionEmailsByUserIDAndStatus(ctx, sub.UserID, domain.SubscriptionStatusUnsubscribed)
			assert.NoError(t, err)
			if tc.isFounded {
				assert.Len(t, result, 1)
			} else {
				assert.Len(t, result, 0)
			}

			suite.rollbackSubscription(t, ctx, sub)
		})

	}
}

func (g *Suite) prepareSubscription(t *testing.T, ctx context.Context, sub domain.Subscription) {
	db := g.db.Model(ctx)
	u := model.User{
		ID:    sub.UserID,
		Email: sub.UserID + "@example.com",
	}
	err := u.Insert(ctx, db, boil.Infer())
	assert.NoError(t, err)

	u = model.User{
		ID:    sub.SubscriberID,
		Email: sub.SubscriberID + "@example.com",
	}
	err = u.Insert(ctx, db, boil.Infer())
	assert.NoError(t, err)
}

func (g *Suite) rollbackSubscription(t *testing.T, ctx context.Context, sub domain.Subscription) {
	db := g.db.Model(ctx)
	if sub.Id != "" {
		sub2 := &model.Subscription{
			ID: sub.Id,
		}
		_, err := sub2.Delete(ctx, db)
		assert.NoError(t, err)
	}

	users := &model.UserSlice{
		{ID: sub.UserID},
		{ID: sub.SubscriberID},
	}
	_, err := users.DeleteAll(ctx, db)
	assert.NoError(t, err)
}
