package middlewares

import (
	"errors"
	"fmt"
	"go_base/utils/log"
	"go_base/utils/rescode"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				//登录认证异常
				if result, ok := r.(*ExceptionResult); ok {
					c.JSON(result.HttpStatus, result.Data)
					c.Abort()
				} else {
					var err error
					var ok bool
					if err, ok = r.(error); !ok {
						err = errors.New(fmt.Sprintf("%s", r))
					}
					//日志记录
					log.Logger.Error(fmt.Sprintf("request url: %s \r %s", c.Request.URL.RequestURI, r))
					//raven.CaptureError(err, map[string]string{"request url": c.Request.RequestURI})

					c.JSON(http.StatusInternalServerError, gin.H{"rescode": rescode.Error, "data": nil, "msg": err.Error()})
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
