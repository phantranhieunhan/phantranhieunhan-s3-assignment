package domain

import "context"

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

func (f FriendshipStatus) CanNotSubscribe() bool {
	switch f {
	case FriendshipStatusBlocked:
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

type Friendships []Friendship

type FriendshipRepo interface {
	Create(ctx context.Context, d Friendship) (string, error)
	UpdateStatus(ctx context.Context, id string, status FriendshipStatus) error
	GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (Friendship, error)
	GetFriendshipByUserIDAndStatus(ctx context.Context, mapEmailUser map[string]string, status ...FriendshipStatus) ([]string, error)
}
