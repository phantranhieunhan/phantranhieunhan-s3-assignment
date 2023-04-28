package friendship

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/adapter/postgres/repository"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/query"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/port"
	graph "github.com/phantranhieunhan/s3-assignment/module/friendship/port/graphql"
)

func New(r *gin.Engine, db postgres.Database) {
	friendshipRepo := repository.NewFriendshipRepository(db)
	userRepo := repository.NewUserRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)

	application := app.Application{
		Commands: app.Commands{
			ConnectFriendship: command.NewConnectFriendshipHandler(friendshipRepo, userRepo, db),
			SubscribeUser:     command.NewSubscribeUserHandler(friendshipRepo, userRepo, subRepo, db),
			BlockUpdatesUser:  command.NewBlockUpdatesUserHandler(friendshipRepo, userRepo, subRepo, db),
		},
		Queries: app.Queries{
			ListFriends:       query.NewListFriendsHandler(friendshipRepo, userRepo),
			ListCommonFriends: query.NewListCommonFriendsHandler(friendshipRepo, userRepo),
			ListUpdatesUser:   query.NewListUpdatesUserHandler(subRepo, userRepo),
			ListSubscriptions: query.NewListSubscriptionsHandler(subRepo),
		},
	}
	port.NewServer(application).Router(r)
	resolver := graph.NewResolver(application)
	handler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))
	r.GET("/", playgroundHandler())
	r.POST("/query", gin.WrapH(handler))
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
