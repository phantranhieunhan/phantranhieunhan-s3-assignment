package port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"
)

type ListUpdatesUserReq struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

func (c ListUpdatesUserReq) validate() error {
	return common.ValidateRequired(c.Sender, "sender")
}

type ListUpdatesUserRes struct {
	Recipients []string `json:"recipients"`
}

func (s *Server) ListUpdatesUser(c *gin.Context) {
	var req ListUpdatesUserReq
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Error("ListUpdatesUser.ShouldBind: ", err)
		common.HttpErrorHandler(c, common.ErrInvalidRequest(err, constant.FRIENDS))
		return
	}

	if err = req.validate(); err != nil {
		logger.Error("ListUpdatesUser.Validate: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	list, err := s.app.Queries.ListUpdatesUser.Handle(c.Request.Context(), req.Sender, req.Text)
	if err != nil {
		logger.Error("ListUpdatesUser.Handle: ", err)
		common.HttpErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.CustomSuccessResponse(
		ListUpdatesUserRes{Recipients: list},
	))
}
