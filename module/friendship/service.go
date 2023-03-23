package friendship

import (
	"github.com/gin-gonic/gin"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/common/adapter/rabbitmq"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/repository"
	friendshiprabbitmq "github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/rabbitmq"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/query"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port/consumer"
)

func New(r *gin.Engine, db postgres.Database, rabbitMQ rabbitmq.MQ) {
	friendshipRepo := repository.NewFriendshipRepository(db)
	userRepo := repository.NewUserRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)

	subMQ := friendshiprabbitmq.NewSubscribeUserMQ(rabbitMQ)
	application := app.Application{
		Commands: app.Commands{
			ConnectFriendship: command.NewConnectFriendshipHandler(friendshipRepo, userRepo, db, subMQ),
			SubscribeUser:     command.NewSubscribeUserHandler(friendshipRepo, userRepo, subRepo, db),
			BlockUpdatesUser:  command.NewBlockUpdatesUserHandler(friendshipRepo, userRepo, subRepo, db),
		},
		Queries: app.Queries{
			ListFriends:       query.NewListFriendsHandler(friendshipRepo, userRepo),
			ListCommonFriends: query.NewListCommonFriendsHandler(friendshipRepo, userRepo),
			ListUpdatesUser:   query.NewListUpdatesUserHandler(subRepo, userRepo),
		},
	}
	consumer := consumer.NewConsumer(rabbitMQ, application)
	consumer.Consume()
	port.NewServer(r, application)
}
