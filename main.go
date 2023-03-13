package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/example/basic/docs"

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

	docs.SwaggerInfo_swagger.Version = "1.0"
	docs.SwaggerInfo_swagger.Title = "EventSourcing Microservice"
	docs.SwaggerInfo_swagger.Description = "EventSourcing CQRS Microservice."
	docs.SwaggerInfo_swagger.Version = "1.0"
	docs.SwaggerInfo_swagger.BasePath = "/api/v1"

	friendship.New(r, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(r.Run(":" + config.C.Server.Port))
}
