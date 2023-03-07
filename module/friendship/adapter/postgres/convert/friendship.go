package convert

import (
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/model"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

func ToFriendshipModel(d domain.Friendship) model.Friendship {
	return model.Friendship{
		Base: model.Base{
			Id: d.Base.Id,
		},
		UserID:   d.UserID,
		FriendID: d.FriendID,
		Status:   d.Status,
	}
}

func ToFriendshipDomain(d model.Friendship) domain.Friendship {
	return domain.Friendship{
		Base: domain.Base{
			Id:        d.Base.Id,
			CreatedAt: d.Base.CreatedAt,
			UpdatedAt: d.Base.UpdatedAt,
		},
		UserID:   d.UserID,
		FriendID: d.FriendID,
		Status:   d.Status,
	}
}
