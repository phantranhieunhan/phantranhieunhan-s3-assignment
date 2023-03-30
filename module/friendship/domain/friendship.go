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
	return f == FriendshipStatusUnfriended
}

func (f FriendshipStatus) CanBlockUser() bool {
	return f == FriendshipStatusUnfriended
}

func (f FriendshipStatus) CanNotSubscribe() bool {
	return f == FriendshipStatusBlocked
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

func (r Friendship) FriendshipWithBlock(userID, friendID string) Friendship {
	return Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   FriendshipStatusBlocked,
	}
}

type Friendships []Friendship

type FriendshipRepo interface {
	Create(ctx context.Context, d Friendship) (string, error)
	UpdateStatus(ctx context.Context, id string, status FriendshipStatus) error
	GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (Friendship, error)
	GetFriendshipByUserIDAndStatus(ctx context.Context, mapEmailUser map[string]string, status ...FriendshipStatus) ([]string, error)
}
