package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ListFriendsRepo interface {
	GetFriendshipByUserIDAndStatus(ctx context.Context, email string, status domain.FriendshipStatus) (domain.Friendships, error)
}

type ListFriendsHandler struct {
	repo ListFriendsRepo
}

func NewListFriendsHandler(repo ListFriendsRepo, transactor Transactor) ListFriendsHandler {
	return ListFriendsHandler{
		repo: repo,
	}
}

// func (h ListFriendsHandler) Handle(ctx context.Context, email string) (string, error) {
// 	// if h.repo.GetFriendshipByUserIDAndStatus(ctx, email, )
// 	// return id, err
// }
