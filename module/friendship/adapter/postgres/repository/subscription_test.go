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
	name                     string
	modifySub                domain.Subscription
	isExisted                bool
	isFounded                bool
	mentionedEmails          []string
	blockedEmails            []string
	isInvalidMentionedEmails bool
	result                   []string
	err                      error
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

			suite.rollbackSubscription(t, ctx, sub, []string{sub.Id})
		})

	}
}

func TestSubscription_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	cs := []SubscriptionTestCase{
		{
			name: "successful",
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

			updateErr := repo.UpdateStatus(ctx, id, domain.SubscriptionStatusUnsubscribed)
			result, err := repo.GetSubscription(ctx, domain.Subscriptions{sub})
			assert.NoError(t, err)
			if tc.err == nil {
				assert.NoError(t, updateErr)
				assert.Equal(t, domain.SubscriptionStatusUnsubscribed, result[0].Status)
			} else {
				assert.Error(t, updateErr)
				assert.Equal(t, domain.SubscriptionStatusSubscribed, result[0].Status)
			}

			suite.rollbackSubscription(t, ctx, sub, []string{sub.Id})
		})
	}
}

func TestSubscription_UpsertSubscription(t *testing.T) {
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
			sub.Id, upsertErr = repo.UpsertSubscription(ctx, upsertSub)
			assert.NoError(t, upsertErr)

			result, err := repo.GetSubscription(ctx, domain.Subscriptions{sub})
			assert.NoError(t, err)
			assert.Equal(t, upsertSub.Status, result[0].Status)

			suite.rollbackSubscription(t, ctx, sub, []string{sub.Id})
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

			suite.rollbackSubscription(t, ctx, sub, []string{sub.Id})
		})

	}
}

func TestGetSubscriptionEmailsByUserIDAndStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewSubscriptionRepository(suite.db)

	cs := []SubscriptionTestCase{
		{
			name:      "successful without mentioned emails",
			isFounded: true,
		},
		{
			name:            "successful with valid mentioned emails",
			isFounded:       true,
			mentionedEmails: []string{"lisa@example.com"},
			result:          []string{"lisa@example.com"},
		},
		{
			name:            "successful with block mentioned emails",
			isFounded:       true,
			mentionedEmails: []string{"lisa@example.com", "andy@example.com"},
			blockedEmails:   []string{"andy@example.com"},
			result:          []string{"lisa@example.com"},
		},
		{
			name:                     "successful with invalid mentioned emails",
			isFounded:                true,
			mentionedEmails:          []string{"lisa@example.com", "andy@example.com"},
			isInvalidMentionedEmails: true, // none existed
			blockedEmails:            []string{"andy@example.com"},
			result:                   []string{},
		},
		{
			name: "not found",
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			emails := []string{}
			if tc.isInvalidMentionedEmails {
				emails = append(emails, util.GenUUID()+"@example.com")
			} else {
				emails = append(emails, tc.mentionedEmails...)
			}

			emails = append(emails, tc.blockedEmails...)
			mapEmailUser := suite.initialUsers(t, ctx, emails)

			subIds := make([]string, 0)
			sub := domain.Subscription{
				UserID:       util.GenUUID(),
				SubscriberID: util.GenUUID(),
				Status:       domain.SubscriptionStatusSubscribed,
			}
			suite.prepareSubscription(t, ctx, sub)
			var err error
			id := "fake-subscription"
			if tc.isFounded {
				id, err = repo.Create(ctx, sub)
				assert.NoError(t, err)
				sub.Id = id
				subIds = append(subIds, id)
			}
			mentionedEmail := make([]string, 0)
			if len(tc.mentionedEmails) > 0 {
				for _, email := range tc.mentionedEmails {
					newEmail := mapEmailUser[email].Email
					if newEmail != "" {
						mentionedEmail = append(mentionedEmail, newEmail)
					}
				}
			}

			if len(tc.blockedEmails) > 0 {
				for _, email := range tc.blockedEmails {
					id, err = repo.UpsertSubscription(ctx, domain.Subscription{
						UserID:       sub.UserID,
						SubscriberID: mapEmailUser[email].ID,
						Status:       domain.SubscriptionStatusUnsubscribed,
					})
					assert.NoError(t, err)
					subIds = append(subIds, id)
				}
			}

			result, err := repo.GetSubscriptionEmailsByUserIDAndEmails(ctx, sub.UserID, mentionedEmail)
			assert.NoError(t, err)
			if tc.isFounded {
				assert.Equal(t, len(tc.result)+1, len(result))
				assert.True(t, util.IsContain(result, sub.SubscriberID+"@example.com"))

				for _, v := range tc.result {
					assert.True(t, util.IsContain(result, mapEmailUser[v].ID+v))
				}
			} else {
				assert.Len(t, result, 0)
			}

			suite.rollbackSubscription(t, ctx, sub, subIds)
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

func (g *Suite) initialUsers(t *testing.T, ctx context.Context, emails []string) map[string]model.User {
	emails = util.RemoveDuplicates(emails)
	result := make(map[string]model.User, 0)
	for _, email := range emails {
		randString := util.GenUUID()
		user := model.User{ID: randString, Email: randString + email}
		err := user.Insert(ctx, g.db.Model(ctx), boil.Infer())
		assert.NoError(t, err)
		result[email] = user
	}
	return result
}

func (g *Suite) rollbackSubscription(t *testing.T, ctx context.Context, sub domain.Subscription, subIds []string) {
	db := g.db.Model(ctx)
	ids := util.RemoveDuplicates(subIds)
	if len(ids) > 0 {
		for _, id := range ids {
			sub2 := &model.Subscription{
				ID: id,
			}
			_, err := sub2.Delete(ctx, db)
			assert.NoError(t, err)
		}
	}

	users := &model.UserSlice{
		{ID: sub.UserID},
		{ID: sub.SubscriberID},
	}
	_, err := users.DeleteAll(ctx, db)
	assert.NoError(t, err)
}
