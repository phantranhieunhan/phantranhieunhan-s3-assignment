package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type FriendshipRepository struct {
	postgres.Database
}

func NewFriendshipRepository(db postgres.Database) FriendshipRepository {
	return FriendshipRepository{
		Database: db,
	}
}

func (f FriendshipRepository) Create(ctx context.Context, d domain.Friendship) (string, error) {
	m := convert.ToFriendshipModel(d)
	m.Base.BeforeCreate()
	if err := f.Model(ctx).Create(&m).Error; err != nil {
		return "", common.ErrDB(err)
	}
	return m.FriendID, nil
}

func (f FriendshipRepository) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	m := convert.ToFriendshipModel(domain.Friendship{
		Base:   domain.Base{Id: id},
		Status: status,
	})
	m.BeforeUpdate()
	if err := f.Model(ctx).Updates(&m).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (f FriendshipRepository) GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error) {
	var m model.Friendship
	if err := f.Model(ctx).Table(model.Friendship{}.TableName()).Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Friendship{}, domain.ErrRecordNotFound
		}
		return domain.Friendship{}, common.ErrDB(err)
	}

	return convert.ToFriendshipDomain(m), nil
}

func (f FriendshipRepository) GetFriendshipByUserIDAndStatus(ctx context.Context, userID string, status ...domain.FriendshipStatus) (domain.Friendships, error) {
	var m model.Friendships
	if err := f.Model(ctx).Table(model.Friendship{}.TableName()).Where("(user_id = ? OR friend_id = ?) AND status IN ?", userID, userID, status).Find(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Friendships{}, domain.ErrRecordNotFound
		}
		return domain.Friendships{}, common.ErrDB(err)
	}

	return convert.ToFriendshipsDomain(m), nil
}
