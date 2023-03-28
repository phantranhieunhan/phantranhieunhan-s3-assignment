package port

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"
)

type ConnectFriendshipReq struct {
	Friends []string `json:"friends"`
}

func (c ConnectFriendshipReq) validate() error {
	if len(c.Friends) != 2 {
		return common.ErrInvalidRequest(fmt.Errorf("friends must be of length 2"), constant.FRIENDS)
	}

	if c.Friends[0] == c.Friends[1] {
		return common.ErrInvalidRequest(fmt.Errorf("friends must be different"), constant.FRIENDS)
	}

	for i, friend := range c.Friends {
		if err := common.ValidateRequired(friend, fmt.Sprintf("friend %d", i)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ConnectFriendship(c *gin.Context) {
	var req ConnectFriendshipReq
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error("ConnectFriendship.ShouldBind: ", err)
		common.HttpErrorHandler(c, common.ErrInvalidRequest(err, constant.FRIENDS))
		return
	}

	if err = req.validate(); err != nil {
		logger.Error("ConnectFriendship.Validate: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	_, err = s.app.Commands.ConnectFriendship.Handle(c.Request.Context(), req.Friends[0], req.Friends[1])
	if err != nil {
		logger.Error("ConnectFriendship.Handle: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(nil))
}
