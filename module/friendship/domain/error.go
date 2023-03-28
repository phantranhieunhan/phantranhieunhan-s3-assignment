package domain

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrUpdateRecordNotFound    = errors.New("update record not found")
	ErrFriendshipIsUnavailable = errors.New("error friendship is unavailable")

	ErrNotFoundUserByEmail = errors.New("not found user by email")

	ErrEmailIsNotValid = errors.New("emails is not valid")
)
