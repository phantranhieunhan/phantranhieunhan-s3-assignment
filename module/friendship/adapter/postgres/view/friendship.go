package view

type Email struct {
	UserEmail   string `json:"user_email" boil:"user_email"`
	FriendEmail string `json:"friend_email" boil:"friend_email"`
}