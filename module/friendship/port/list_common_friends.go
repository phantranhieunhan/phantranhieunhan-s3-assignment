package port

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"
)

type ListCommonFriendsReq struct {
	Friends []string `json:"friends"`
}

func (l ListCommonFriendsReq) validate() error {
	if len(l.Friends) != 2 {
		return common.ErrInvalidRequest(fmt.Errorf("friends must be of length 2"), constant.FRIENDS)
	}

	for i, friend := range l.Friends {
		if err := common.ValidateRequired(friend, fmt.Sprintf("friend %d", i)); err != nil {
			return err
		}
		if err := common.ValidateEmail(friend); err != nil {
			return err
		}
	}
	return nil
}

type ListCommonFriendsResp struct {
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (s *Server) ListCommonFriends(c *gin.Context) {
	var req ListCommonFriendsReq
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error("ListCommonFriends.ShouldBind: ", err)
		common.HttpErrorHandler(c, common.ErrInvalidRequest(err, constant.FRIENDS))
		return
	}

	if err = req.validate(); err != nil {
		logger.Error("ListCommonFriends.Validate: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	list, err := s.app.Queries.ListCommonFriends.Handle(c.Request.Context(), req.Friends)
	if err != nil {
		logger.Error("ListFriends.Handle: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.CustomSuccessResponse(
		ListCommonFriendsResp{
			Friends: list,
			Count:   len(list),
		},
	))
}
