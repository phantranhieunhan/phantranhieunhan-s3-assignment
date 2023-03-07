package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ConnectFriendshipRepo interface {
	Create(ctx context.Context, d domain.Friendship) (string, error)
	GetFriendshipByUserID(ctx context.Context, userID, friendID string) (domain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error
}

type ConnectFriendshipHandler struct {
	repo       ConnectFriendshipRepo
	transactor Transactor
}

func NewConnectFriendshipHandler(repo ConnectFriendshipRepo, transactor Transactor) ConnectFriendshipHandler {
	return ConnectFriendshipHandler{
		repo:       repo,
		transactor: transactor,
	}
}

func (h ConnectFriendshipHandler) Create(ctx context.Context, d domain.Friendship) (string, error) {
	var id string
	d.Status = domain.FriendshipStatusFriended
	err := h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := h.repo.GetFriendshipByUserID(ctx, d.UserID, d.FriendID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserID %w", err)
			return common.ErrCannotGetEntity(d.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound {
			id, err = h.repo.Create(ctx, d)
			if err != nil {
				logger.Errorf("repo.Create %w", err)
				return common.ErrCannotCreateEntity(d.DomainName(), err)
			}
		} else {
			if !f.Status.CanConnect() {
				logger.Errorf("Status.CanConnect")
				return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
			}
			if err = h.repo.UpdateStatus(ctx, f.Id, domain.FriendshipStatusFriended); err != nil {
				logger.Errorf("repo.UpdateStatus %w", err)
				return common.ErrCannotUpdateEntity(d.DomainName(), err)
			}
			id = f.Id
		}
		return err
	})

	return id, err
}
