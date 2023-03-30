package convert

import (
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

func ToFriendshipModel(d domain.Friendship) model.Friendship {
	return model.Friendship{
		ID:       d.Id,
		UserID:   d.UserID,
		FriendID: d.FriendID,
		Status:   int(d.Status),
	}
}

func ToFriendshipDomain(d model.Friendship) domain.Friendship {
	return domain.Friendship{
		Base: domain.Base{
			Id:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		UserID:   d.UserID,
		FriendID: d.FriendID,
		Status:   domain.FriendshipStatus(d.Status),
	}
}

func ToFriendshipsDomain(m model.FriendshipSlice) domain.Friendships {
	ds := make(domain.Friendships, 0, len(m))
	for _, v := range m {
		ds = append(ds, ToFriendshipDomain(*v))
	}
	return ds
}

func ToMapUserEmailDomainList(users model.UserSlice) map[string]string {
	result := make(map[string]string, 0)
	for _, v := range users {
		result[v.ID] = v.Email
	}
	return result
}

func ToMapEmailUserDomainList(users model.UserSlice) map[string]string {
	result := make(map[string]string, 0)
	for _, v := range users {
		result[v.Email] = v.ID
	}
	return result
}
