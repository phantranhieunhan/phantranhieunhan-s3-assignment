package port

import (
	"net/http"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"

	"github.com/gin-gonic/gin"
)

type BlockUpdatesUserReq struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

func (l BlockUpdatesUserReq) validate() error {
	if err := common.ValidateRequired(l.Requestor, constant.REQUESTOR); err != nil {
		return err
	}
	if err := common.ValidateEmail(l.Requestor); err != nil {
		return err
	}

	if err := common.ValidateRequired(l.Target, constant.TARGET); err != nil {
		return err
	}
	if err := common.ValidateEmail(l.Target); err != nil {
		return err
	}

	return nil
}

func (s *Server) BlockUpdatesUser(c *gin.Context) {
	var req BlockUpdatesUserReq
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error("ListCommonFriends.ShouldBind: ", err)
		common.HttpErrorHandler(c, common.ErrInvalidRequest(err, "body data"))
		return
	}

	if err = req.validate(); err != nil {
		logger.Error("ListCommonFriends.Validate: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	err = s.app.Commands.BlockUpdatesUser.Handle(c.Request.Context(), command.BlockUpdatesUserPayload{
		Requestor: req.Requestor,
		Target:    req.Target,
	})
	if err != nil {
		logger.Error("BlockUpdatesUser.Handle: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.CustomSuccessResponse(nil))
}
