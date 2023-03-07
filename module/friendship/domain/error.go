package domain

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrFriendshipIsUnavailable = errors.New("error friendship is unavailable")
)
