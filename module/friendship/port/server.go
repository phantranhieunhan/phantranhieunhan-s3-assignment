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
	friendship := r.Group("friendship")
	friendship.POST("connect", s.ConnectFriendship)
	friendship.GET("friends", s.ListFriends)
	friendship.GET("mutuals", s.ListCommonFriends)

	subscription := r.Group("subscription")
	subscription.POST("subscribe", s.SubscribeUser)
	subscription.POST("block", s.BlockUpdatesUser)
	return s
}
