package domain

type Follower struct {
	Base        Base   `json:",inline"`
	UserID      string `json:"user_id"`
	FollowingID string `json:"following_id"`
}
