package domain

import "context"

type User struct {
	Base     Base   `json:",inline"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r User) DomainName() string {
	return "User"
}

type UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}
