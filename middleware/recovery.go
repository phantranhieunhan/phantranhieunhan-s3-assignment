package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/common/constant"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Header("Content-Type", "application/json")

			if appErr, ok := err.(*common.AppError); ok {
				if config.C.Env == constant.PRODUCTION_ENV_NAME {
					appErr.ClearRoot()
				}
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				return
			}

			appErr := common.ErrInternal(err.(error))
			c.AbortWithStatusJSON(appErr.StatusCode, appErr)
			panic(err)
			// return
		}
	}()
	c.Next()
}
