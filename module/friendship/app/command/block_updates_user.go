package command

import (
	"context"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/domain"
)

type BlockUpdatesUserPayload struct {
	Requestor string
	Target    string
}

type BlockUpdatesUserHandler struct {
	friendshipRepo FriendshipRepo
	userRepo       UserRepo
	transactor     Transactor
}

func (b BlockUpdatesUserHandler) Handle(ctx context.Context, payload BlockUpdatesUserPayload) error {
	if payload.Requestor == payload.Target {
		return common.ErrInvalidRequest(domain.ErrEmailIsNotValid, "payload")
	}

	userIDs, _, err := b.userRepo.GetUserIDsByEmails(ctx, []string{payload.Requestor, payload.Target})
	if err != nil {
		logger.Errorf("userRepo.GetUserIDsByEmails %w", err)
		if err == domain.ErrNotFoundUserByEmail {
			return common.ErrInvalidRequest(err, "emails")
		}
		return common.ErrCannotGetEntity(domain.Subscription{}.DomainName(), err)
	}

	err = b.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		f, err := b.friendshipRepo.GetFriendshipByUserIDs(ctx, userIDs[payload.Requestor], userIDs[payload.Target])
		if err != nil && err != domain.ErrRecordNotFound {
			logger.Errorf("Create.GetFriendshipByUserIDs %w", err)
			return common.ErrCannotGetEntity(f.DomainName(), err)
		}

		if err == domain.ErrRecordNotFound {
			id, err = b.friendshipRepo.Create(ctx, d)
			if err != nil {
				logger.Errorf("repo.Create %w", err)
				return common.ErrCannotCreateEntity(d.DomainName(), err)
			}
		} else {
			if !f.Status.CanConnect() {
				logger.Errorf("Status.CanConnect")
				return common.ErrInvalidRequest(domain.ErrFriendshipIsUnavailable, "")
			}
			if err = b.friendshipRepo.UpdateStatus(ctx, f.Id, domain.FriendshipStatusFriended); err != nil {
				logger.Errorf("repo.UpdateStatus %w", err)
				return common.ErrCannotUpdateEntity(d.DomainName(), err)
			}
			id = f.Id
		}
		return nil
	})
}
