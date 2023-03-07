package domain

type FriendshipStatus int

const (
	FriendshipStatusInvalid FriendshipStatus = iota
	FriendshipStatusFriended
	FriendshipStatusPending
	FriendshipStatusUnfriended
	FriendshipStatusBlocked
)

func (f FriendshipStatus) CanConnect() bool {
	switch f {
	case FriendshipStatusUnfriended:
		return true
	default:
		return false
	}
}

type Friendship struct {
	Base     `json:",inline"`
	UserID   string           `json:"user_id"`
	FriendID string           `json:"friend_id"`
	Status   FriendshipStatus `json:"status"`
}

func (r Friendship) DomainName() string {
	return "Friendship"
}
