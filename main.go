package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phantranhieunhan/s3-assignment/common/adapter/postgres"
	"github.com/phantranhieunhan/s3-assignment/common/logger"
	"github.com/phantranhieunhan/s3-assignment/middleware"
	"github.com/phantranhieunhan/s3-assignment/module/friendship"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
)

func main() {
	flag.Parse()
	err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
	}

	// Init logger.
	logger.Setup(config.C.Env)

	db := postgres.NewDatabase()

	r := gin.Default()
	r.Use(middleware.Recover)

	friendship.New(r, db)

	log.Fatal(r.Run(":" + config.C.Server.Port))
}
