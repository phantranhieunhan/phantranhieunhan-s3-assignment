package repository

import (
	"context"
	"database/sql"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserRepository struct {
	db postgres.Database
}

func NewUserRepository(db postgres.Database) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (f UserRepository) GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, []string, error) {
	users, err := model.Users(AndIn("email IN ?", util.InterfaceSlice(emails)...)).All(ctx, f.db.Model(ctx))
	// users, err := model.Users().All(ctx, f.db.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, domain.ErrRecordNotFound
		}
		return nil, nil, common.ErrDB(err)
	}
	if len(users) != len(emails) {
		return nil, nil, domain.ErrNotFoundUserByEmail
	}
	result := make(map[string]string, 0)
	sliceString := make([]string, 0, len(users))
	for _, v := range users {
		result[v.Email] = v.ID
		sliceString = append(sliceString, v.ID)
	}

	return result, sliceString, nil
}

func (f UserRepository) GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, []string, error) {
	users, err := model.Users(AndIn("id IN ?", util.InterfaceSlice(userIDs)...)).All(ctx, f.db.Model(ctx))
	if err != nil {
		return nil, nil, common.ErrDB(err)
	}

	if len(users) != len(userIDs) {
		return nil, nil, domain.ErrNotFoundUserByEmail
	}

	result := make(map[string]string, 0)
	emails := make([]string, 0, len(users))
	for _, v := range users {
		result[v.ID] = v.Email
		emails = append(emails, v.Email)
	}

	return result, emails, nil
}
