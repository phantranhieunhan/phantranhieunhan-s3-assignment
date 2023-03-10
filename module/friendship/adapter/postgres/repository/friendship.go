package repository

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
	"github.com/phantranhieunhan/s3-assignment/pkg/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type FriendshipRepository struct {
	db postgres.Database
}

func NewFriendshipRepository(db postgres.Database) FriendshipRepository {
	return FriendshipRepository{
		db: db,
	}
}

func (f FriendshipRepository) Create(ctx context.Context, d domain.Friendship) (string, error) {
	d.Id = util.GenUUID()
	m := convert.ToFriendshipModel(d)
	if err := m.Insert(ctx, f.db.DB, boil.Infer()); err != nil {
		return "", common.ErrDB(err)
	}
	return m.FriendID, nil
}

func (f FriendshipRepository) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	m := convert.ToFriendshipModel(domain.Friendship{
		Base:   domain.Base{Id: id},
		Status: status,
	})
	_, err := m.Update(ctx, f.db.DB, boil.Whitelist(model.FollowerColumns.Status, model.FollowerColumns.UpdatedAt))
	if err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (f FriendshipRepository) GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error) {
	m, err := model.Friendships(qm.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID)).All(ctx, f.db.DB)

	if err != nil {
		return domain.Friendship{}, common.ErrDB(err)
	}
	if len(m) == 0 {
		return domain.Friendship{}, domain.ErrRecordNotFound
	}
	return convert.ToFriendshipDomain(*(m[0])), nil
}

func (f FriendshipRepository) GetFriendshipByUserIDAndStatus(ctx context.Context, userID string, status ...domain.FriendshipStatus) (domain.Friendships, error) {
	m, err := model.Friendships(qm.Where("(user_id = ? OR friend_id = ?)", userID, userID), qm.AndIn("status IN ?", util.InterfaceSlice(status)...)).All(ctx, f.db.DB)
	if err != nil {
		return domain.Friendships{}, common.ErrDB(err)
	}

	return convert.ToFriendshipsDomain(m), nil
}
