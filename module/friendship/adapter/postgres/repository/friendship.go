package repository

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/convert"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/view"
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
	if err := m.Insert(ctx, f.db.Model(ctx), boil.Infer()); err != nil {
		return "", common.ErrDB(err)
	}
	return m.ID, nil
}

func (f FriendshipRepository) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	m := model.Friendship{
		ID:     id,
		Status: int(status),
	}
	_, err := m.Update(ctx, f.db.Model(ctx), boil.Whitelist(model.FriendshipColumns.Status, model.FriendshipColumns.UpdatedAt))
	if err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (f FriendshipRepository) GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error) {
	m, err := model.Friendships(qm.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID)).All(ctx, f.db.Model(ctx))

	if err != nil {
		return domain.Friendship{}, common.ErrDB(err)
	}
	if len(m) == 0 {
		return domain.Friendship{}, domain.ErrRecordNotFound
	}
	return convert.ToFriendshipDomain(*(m[0])), nil
}

func (f FriendshipRepository) GetFriendshipByUserIDAndStatus(ctx context.Context, mapEmailUser map[string]string, status ...domain.FriendshipStatus) ([]string, error) {

	var (
		resultEmails []view.Email
		emptyList    = []string{}
	)

	where := []qm.QueryMod{
		qm.Select("u1.email as user_email", "u2.email as friend_email"),
		qm.From("friendships f"),
		qm.LeftOuterJoin("users u1 on f.user_id = u1.id"),
		qm.LeftOuterJoin("users u2 on f.friend_id = u2.id"),
	}
	userIDs := util.MapValuesToSlice(mapEmailUser)
	for _, userID := range userIDs {
		where = append(where, qm.Or("(f.user_id = ? OR f.friend_id = ?)", userID, userID))
	}
	statusList, err := util.InterfaceSlice(status)
	if err != nil {
		return emptyList, err
	}

	where = append(where, qm.AndIn("status IN ?", statusList...))

	err = model.NewQuery(where...).Bind(ctx, f.db.Model(ctx), &resultEmails)
	if err != nil {
		return emptyList, common.ErrDB(err)
	}

	if len(resultEmails) == 0 {
		return emptyList, domain.ErrRecordNotFound
	}

	result := f.getFriendsOfUsers(resultEmails, mapEmailUser)

	return result, nil
}

// get friendIDs from userId or friendId field if not same userID
func (f FriendshipRepository) getFriendsOfUsers(entireEmails []view.Email, mapEmailUser map[string]string) []string {
	result := make([]string, 0)
	emails := util.MapKeysToSlice(mapEmailUser)
	for _, v := range entireEmails {
		var y string
		if util.IsContain(emails, v.UserEmail) {
			y = v.FriendEmail
		}
		if util.IsContain(emails, v.FriendEmail) {
			y = v.UserEmail
		}
		result = append(result, y)
	}
	return result
}
