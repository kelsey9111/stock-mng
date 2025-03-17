package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"stock-management/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.Error("Recovered from panic",
					zap.String("error", fmt.Sprintf("%v", err)),
					zap.String("stack", string(debug.Stack())),
				)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
