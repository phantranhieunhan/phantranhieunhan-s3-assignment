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

func (f UserRepository) GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error) {
	users, err := model.Users(AndIn("email IN ?", util.InterfaceSlice(emails)...)).All(ctx, f.db.DB)
	// users, err := model.Users().All(ctx, f.db.DB)
	if err != nil {
		zResult := make(map[string]string, 0)
		if err == sql.ErrNoRows {
			return zResult, domain.ErrRecordNotFound
		}
		return zResult, common.ErrDB(err)
	}
	if len(users) != len(emails) {
		return nil, domain.ErrNotFoundUserByEmail
	}
	result := make(map[string]string, 0)
	for _, v := range users {
		result[v.Email] = v.ID
	}

	return result, nil
}

func (f UserRepository) GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	users, err := model.Users(AndIn("id IN ?", util.InterfaceSlice(userIDs)...)).All(ctx, f.db.DB)
	zResult := make(map[string]string, 0)
	if err != nil {
		return zResult, common.ErrDB(err)
	}
	l := len(users)
	if l == 0 || l != len(userIDs) {
		return nil, domain.ErrNotFoundUserByEmail
	}

	result := make(map[string]string, 0)
	for _, v := range users {
		result[v.ID] = v.Email
	}

	return result, nil
}
