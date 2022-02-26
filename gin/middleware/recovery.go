package middleware

import (
	"net/http"

	"e.coding.net/guodf/gopkg/goutil/gin/dtos"
	"github.com/gin-gonic/gin"
)

// Recovery 全局错误中间件
func Recovery() gin.RecoveryFunc {
	return func(c *gin.Context, err interface{}) {
		if code, ok := err.(int); ok {
			c.AbortWithStatusJSON(http.StatusOK, dtos.FailedResult(dtos.ErrorCode(code)))
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
