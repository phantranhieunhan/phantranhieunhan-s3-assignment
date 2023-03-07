package model

import "github.com/phantranhieunhan/s3-assignment/module/friendship/domain"

type Friendship struct {
	Base     `json:",inline"`
	UserID   string                  `json:"user_id" gorm:"column:user_id"`
	FriendID string                  `json:"friend_id" gorm:"column:friend_id"`
	Status   domain.FriendshipStatus `json:"status" gorm:"column:status"`
}

func (r Friendship) TableName() string {
	return "friendships"
}
