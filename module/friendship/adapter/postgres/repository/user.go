package repository

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	iEmails, err := util.InterfaceSlice(emails)
	if err != nil {
		return nil, common.ErrInvalidRequest(err, "userIDs")
	}
	users, err := model.Users(qm.AndIn("email IN ?", iEmails...)).All(ctx, f.db.Model(ctx))
	if err != nil {
		return nil, common.ErrDB(err)
	}
	if len(users) != len(emails) {
		return nil, domain.ErrNotFoundUserByEmail
	}
	result := convert.ToMapEmailUserDomainList(users)

	return result, nil
}

func (f UserRepository) GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	iUserIDs, err := util.InterfaceSlice(userIDs)
	emptyResult := make(map[string]string, 0)
	if err != nil {
		return emptyResult, common.ErrInvalidRequest(err, "userIDs")
	}
	users, err := model.Users(qm.AndIn("id IN ?", iUserIDs...)).All(ctx, f.db.Model(ctx))

	if err != nil {
		return emptyResult, common.ErrDB(err)
	}

	if len(users) != len(userIDs) {
		return nil, domain.ErrNotFoundUserByEmail
	}
	result := convert.ToMapUserEmailDomainList(users)

	return result, nil
}
