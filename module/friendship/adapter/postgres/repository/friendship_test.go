package repository

import (
	"context"
	"testing"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type FriendshipTestCase struct {
	name                    string
	modifyFriendship        domain.Friendship
	IsErrExistingFriendship bool
	hasError                bool
}

func TestFriendship_Create(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewFriendshipRepository(suite.db)

	cs := []FriendshipTestCase{
		{
			name: "successful",
		},
		{
			name: "fail by reference user id",
			modifyFriendship: domain.Friendship{
				UserID: "user-id",
			},
			hasError: true,
		},
		{
			name: "fail by reference subscriber id",
			modifyFriendship: domain.Friendship{
				FriendID: "friend-id",
			},
			hasError: true,
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			fri := domain.Friendship{
				UserID:   util.GenUUID(),
				FriendID: util.GenUUID(),
				Status:   domain.FriendshipStatusFriended,
			}
			suite.prepareFriendship(t, ctx, fri)
			created := fri
			if tc.hasError {
				if tc.modifyFriendship.FriendID != "" {
					created.FriendID = "friend-id"
				} else if tc.modifyFriendship.UserID != "" {
					created.UserID = "user-id"
				}
			}
			var err error

			fri.Id, err = repo.Create(ctx, created)
			if tc.hasError {
				assert.Error(t, err)
				assert.Empty(t, fri.Id)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, fri.Id)
			}

			suite.rollbackFriendship(t, ctx, fri)
		})

	}
}

func TestFriendship_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewFriendshipRepository(suite.db)

	cs := []FriendshipTestCase{
		{
			name: "successful",
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			fri := domain.Friendship{
				UserID:   util.GenUUID(),
				FriendID: util.GenUUID(),
				Status:   domain.FriendshipStatusFriended,
			}
			suite.prepareFriendship(t, ctx, fri)
			var err error
			fri.Id, err = repo.Create(ctx, fri)
			assert.NoError(t, err)
			id := fri.Id
			if tc.modifyFriendship.Id != "" {
				id = tc.modifyFriendship.Id
			}

			updatedErr := repo.UpdateStatus(ctx, id, domain.FriendshipStatusBlocked)
			result, err := repo.GetFriendshipByUserIDs(ctx, fri.FriendID, fri.UserID)
			assert.NoError(t, err)
			if tc.hasError {
				assert.Error(t, updatedErr)
				assert.Equal(t, domain.FriendshipStatusFriended, result.Status)
			} else {
				assert.NoError(t, updatedErr)
				assert.Equal(t, domain.FriendshipStatusBlocked, result.Status)
			}

			suite.rollbackFriendship(t, ctx, fri)
		})
	}
}

func TestGetFriendshipByUserIDAndStatus(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewFriendshipRepository(suite.db)

	cs := []FriendshipTestCase{
		{
			name: "successful",
		},
		{
			name:     "fail by invalid id",
			hasError: true,
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			fri := domain.Friendship{
				UserID:   util.GenUUID(),
				FriendID: util.GenUUID(),
				Status:   domain.FriendshipStatusFriended,
			}
			suite.prepareFriendship(t, ctx, fri)
			var err error
			if !tc.hasError {
				fri.Id, err = repo.Create(ctx, fri)
				assert.NoError(t, err)
			}

			mapEmailUser := map[string]string{
				fri.UserID + "@example.com": fri.UserID,
			}
			result, err := repo.GetFriendshipByUserIDAndStatus(ctx, mapEmailUser, domain.FriendshipStatusFriended)
			if tc.hasError {
				assert.Len(t, result, 0)
				assert.Equal(t, err, domain.ErrRecordNotFound)
			} else {
				assert.Equal(t, result[0], fri.FriendID+"@example.com")
				assert.NoError(t, err)
			}

			suite.rollbackFriendship(t, ctx, fri)
		})
	}
}

func (g *Suite) prepareFriendship(t *testing.T, ctx context.Context, sub domain.Friendship) {
	db := g.db.Model(ctx)
	u := model.User{
		ID:    sub.UserID,
		Email: sub.UserID + "@example.com",
	}
	err := u.Insert(ctx, db, boil.Infer())
	assert.NoError(t, err)

	u = model.User{
		ID:    sub.FriendID,
		Email: sub.FriendID + "@example.com",
	}
	err = u.Insert(ctx, db, boil.Infer())
	assert.NoError(t, err)
}

func (g *Suite) rollbackFriendship(t *testing.T, ctx context.Context, sub domain.Friendship) {
	db := g.db.Model(ctx)
	if sub.Id != "" {
		sub2 := &model.Friendship{
			ID: sub.Id,
		}
		_, err := sub2.Delete(ctx, db)
		assert.NoError(t, err)
	}

	users := &model.UserSlice{
		{ID: sub.UserID},
		{ID: sub.FriendID},
	}
	_, err := users.DeleteAll(ctx, db)
	assert.NoError(t, err)
}
