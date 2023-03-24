package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			common.HttpErrorHandler(c, err)
			panic(err)
		}
	}()
	c.Next()
}
