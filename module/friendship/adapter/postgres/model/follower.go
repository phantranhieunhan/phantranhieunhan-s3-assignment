package model

type Follower struct {
	Base        `json:",inline"`
	UserID      string `json:"user_id" gorm:"column:user_id"`
	FollowingID string `json:"following_id" gorm:"column:following_id"`
}

func (r Follower) TableName() string {
	return "followers"
}
