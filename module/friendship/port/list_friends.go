package port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"
)

type ListFriendsReq struct {
	Email string `json:"email"`
}

func (c ListFriendsReq) validate() error {
	if err := common.ValidateRequired(c.Email, "email"); err != nil {
		return err
	}

	if err := common.ValidateEmail(c.Email); err != nil {
		return err
	}

	return nil
}

type ListFriendsRes struct {
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (s *Server) ListFriends(c *gin.Context) {
	var req ListFriendsReq
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error("ListFriends.ShouldBind: ", err)
		common.HttpErrorHandler(c, common.ErrInvalidRequest(err, constant.FRIENDS))
		return
	}

	if err = req.validate(); err != nil {
		logger.Error("ListFriends.Validate: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	list, err := s.app.Queries.ListFriends.Handle(c.Request.Context(), req.Email)
	if err != nil {
		logger.Error("ListFriends.Handle: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.CustomSuccessResponse(
		ListFriendsRes{Friends: list, Count: len(list)},
	))
}
