package repository

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	postgres.Database
}

func NewUserRepository(db postgres.Database) UserRepository {
	return UserRepository{
		Database: db,
	}
}

func (f UserRepository) GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error) {
	var users model.Users
	if err := f.Model(ctx).Table(model.User{}.TableName()).Where("email IN ?", emails).Find(&users).Error; err != nil {
		zResult := make(map[string]string, 0)
		if err == gorm.ErrRecordNotFound {
			return zResult, domain.ErrRecordNotFound
		}
		return zResult, common.ErrDB(err)
	}
	if len(users) != len(emails) {
		return nil, domain.ErrNotFoundUserByEmail
	}
	result := make(map[string]string, 0)
	for _, v := range users {
		result[v.Email] = v.Base.Id
	}

	return result, nil
}
