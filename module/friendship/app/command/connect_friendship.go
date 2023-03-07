package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ConnectFriendshipRepo interface {
	Create(ctx context.Context, d domain.Friendship) (string, error)
	GetFriendshipByUserIDs(ctx context.Context, userID, friendID string) (domain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error
}

type UserRepo interface {
	GetUserIDsByEmails(ctx context.Context, emails []string) (map[string]string, error)
}

type ConnectFriendshipHandler struct {
	friendshipRepo ConnectFriendshipRepo
	userRepo       UserRepo
	transactor     Transactor
}

func NewConnectFriendshipHandler(repo ConnectFriendshipRepo, userRepo UserRepo, transactor Transactor) ConnectFriendshipHandler {
	return ConnectFriendshipHandler{
		friendshipRepo: repo,
		userRepo:       userRepo,
		transactor:     transactor,
	}
}

func (h ConnectFriendshipHandler) Handle(ctx context.Context, userEmail, friendEmail string) (string, error) {
	userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, []string{userEmail, friendEmail})
	if err != nil {
		return "", common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}
	var id string
	d := domain.Friendship{
		Status: domain.FriendshipStatusFriended,
		UserID: userIDs[userEmail],
		FriendID: userIDs[friendEmail],
	}

	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := h.friendshipRepo.GetFriendshipByUserIDs(ctx, d.UserID, d.FriendID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(d.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound {
			id, err = h.friendshipRepo.Create(ctx, d)
			if err != nil {
				logger.Errorf("repo.Create %w", err)
				return common.ErrCannotCreateEntity(d.DomainName(), err)
			}
		} else {
			if !f.Status.CanConnect() {
				logger.Errorf("Status.CanConnect")
				return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
			}
			if err = h.friendshipRepo.UpdateStatus(ctx, f.Id, domain.FriendshipStatusFriended); err != nil {
				logger.Errorf("repo.UpdateStatus %w", err)
				return common.ErrCannotUpdateEntity(d.DomainName(), err)
			}
			id = f.Id
		}
		return err
	})

	return id, err
}
