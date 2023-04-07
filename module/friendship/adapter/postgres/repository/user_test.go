package repository

import (
	"context"
	"testing"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/stretchr/testify/assert"
)

type UserTestCase struct {
	name string
	err  error
}

func TestGetEmailsByUserIDs(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewUserRepository(suite.db)

	cs := []UserTestCase{
		{
			name: "successful",
		},
		{
			name: "not found",
			err:  domain.ErrNotFoundUserByEmail,
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			var ids []string
			if tc.err == nil {
				ids = []string{"cd2543cd-6566-4661-a122-2c963fc16b7c", "afed6e29-07d1-443a-a0c7-38d77ef8f332"}
			} else {
				ids = []string{"fake-id"}
			}

			mapUserEmails, err := repo.GetEmailsByUserIDs(ctx, ids)
			if tc.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, "andy@example.com", mapUserEmails["cd2543cd-6566-4661-a122-2c963fc16b7c"])
				assert.Equal(t, "lisa@example.com", mapUserEmails["afed6e29-07d1-443a-a0c7-38d77ef8f332"])
			} else {
				assert.Error(t, err)
			}
		})

	}
}

func TestGetUserIDsByEmails(t *testing.T) {
	ctx := context.Background()
	suite := NewSuite(ctx)
	repo := NewUserRepository(suite.db)

	cs := []UserTestCase{
		{
			name: "successful",
		},
		{
			name: "not found",
			err:  domain.ErrNotFoundUserByEmail,
		},
	}
	for _, tc := range cs {
		t.Run(tc.name, func(t *testing.T) {
			var emails []string
			if tc.err == nil {
				emails = []string{"andy@example.com", "lisa@example.com"}
			} else {
				emails = []string{"fake-email"}
			}

			mapEmailUsers, err := repo.GetUserIDsByEmails(ctx, emails)
			if tc.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, "cd2543cd-6566-4661-a122-2c963fc16b7c", mapEmailUsers["andy@example.com"])
				assert.Equal(t, "afed6e29-07d1-443a-a0c7-38d77ef8f332", mapEmailUsers["lisa@example.com"])
			} else {
				assert.Error(t, err)
			}
		})

	}
}
