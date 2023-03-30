package common

import (
	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common/constant"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
)

func HttpErrorHandler(c *gin.Context, err interface{}) {
	c.Header("Content-Type", "application/json")

	if appErr, ok := err.(*AppError); ok {
		if config.C.Env == constant.PRODUCTION_ENV_NAME {
			appErr.ClearRoot()
		}
		c.AbortWithStatusJSON(appErr.StatusCode, appErr)
		return
	}

	appErr := ErrInternal(err.(error))
	c.AbortWithStatusJSON(appErr.StatusCode, appErr)
}
