package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type ConnectFriendshipHandler struct {
	friendshipRepo domain.FriendshipRepo
	userRepo       domain.UserRepo
	transactor     Transactor
}

func NewConnectFriendshipHandler(repo domain.FriendshipRepo, userRepo domain.UserRepo, transactor Transactor) ConnectFriendshipHandler {
	return ConnectFriendshipHandler{
		friendshipRepo: repo,
		userRepo:       userRepo,
		transactor:     transactor,
	}
}

func (h ConnectFriendshipHandler) Handle(ctx context.Context, userEmail, friendEmail string) (domain.Friendship, error) {
	userIDs, err := h.userRepo.GetUserIDsByEmails(ctx, []string{userEmail, friendEmail})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return domain.Friendship{}, common.ErrInvalidRequest(err, "emails")
		}
		return domain.Friendship{}, common.ErrCannotGetEntity(domain.User{}.DomainName(), err)
	}
	d := domain.Friendship{
		Status:   domain.FriendshipStatusFriended,
		UserID:   userIDs[userEmail],
		FriendID: userIDs[friendEmail],
	}

	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := h.friendshipRepo.GetFriendshipByUserIDs(ctx, d.UserID, d.FriendID)
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(d.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound {
			d.Id, err = h.friendshipRepo.Create(ctx, d)
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
			d.Id = f.Id
		}
		return nil
	})

	if err != nil {
		return domain.Friendship{}, err
	}

	return d, err
}
