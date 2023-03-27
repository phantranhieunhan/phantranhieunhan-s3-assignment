package repository

import (
	"context"
	"log"
	"testing"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
	}
	db := postgres.NewDatabase()
	repo := NewSubscriptionRepository(db)
	ctx := context.Background()

	user := domain.Subscription{
		UserID:       util.GenUUID(),
		SubscriberID: util.GenUUID(),
	}

	util.LoadTestSQLFile(t, db.Model(ctx), "./sql_test/create_subscription_up.sql", user.UserID, user.UserID+"@example.com", user.SubscriberID, user.SubscriberID+"@example.com")

	sub := domain.Subscription{
		UserID:       user.UserID,
		SubscriberID: user.SubscriberID,
		Status:       domain.SubscriptionStatusSubscribed,
	}
	sub.Id, err = repo.Create(ctx, sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, sub.Id)
	
	util.LoadTestSQLFile(t, db.Model(ctx), "./sql_test/create_subscription_dow.sql", sub.Id, user.UserID, user.SubscriberID)
}
