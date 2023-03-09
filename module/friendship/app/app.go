package app

import "github.com/phantranhieunhan/s3-assignment/module/friendship/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ConnectFriendship command.ConnectFriendshipHandler
}

type Queries struct {
	ListFriends command.ListFriendsHandler
}
