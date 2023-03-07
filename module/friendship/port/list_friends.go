package port

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/phantranhieunhan/s3-assignment/common"
// 	"github.com/phantranhieunhan/s3-assignment/common/logger"
// 	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/constant"
// )

// type ListFriendsReq struct {
// 	Email string `json:"email"`
// }

// func (c ListFriendsReq) validate() error {
// 	return common.ValidateRequired(c.Email, "email")
// }

// func (s *Server) ListFriends(c *gin.Context) {
// 	var req ListFriendsReq
// 	var err error
// 	if err = c.ShouldBind(&req); err != nil {
// 		logger.Error("ListFriends.ShouldBind: ", err)
// 		panic(common.ErrInvalidRequest(err, constant.FRIENDS))
// 	}

// 	if err = req.validate(); err != nil {
// 		logger.Error("ListFriends.Validate: ", err)
// 		panic(err)
// 	}

// 	_, err = s.app.Queries.ListFriends.Handle(c.Request.Context(), req.Email)
// 	if err != nil {
// 		logger.Error("ListFriends.Create: ", err)
// 		panic(err)
// 	}

// 	c.JSON(http.StatusOK, common.SimpleSuccessResponse(nil))
// }
