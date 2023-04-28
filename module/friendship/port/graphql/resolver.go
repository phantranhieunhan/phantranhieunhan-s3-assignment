package graph

import "github.com/phantranhieunhan/s3-assignment/module/friendship/app"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	app app.Application
}

func NewResolver(app app.Application) Resolver {
	s := Resolver{app: app}

	return s
}
