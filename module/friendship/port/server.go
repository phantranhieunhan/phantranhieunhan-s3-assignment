package port

import (
	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
)

type Server struct {
	app app.Application
}

func NewServer(r *gin.Engine, app app.Application) Server {
	s := Server{app: app}
	g := r.Group("friendship")
	g.POST("connect", s.ConnectFriendship)
	return s
}
