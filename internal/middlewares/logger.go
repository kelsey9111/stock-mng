package middlewares

import (
	"bytes"
	"context"
	"io"
	"stock-management/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader("X-Trace-ID")
		if traceId == "" {
			traceId = uuid.NewString()
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "TraceID", traceId)

		logger := global.Logger.AddTraceID(ctx)
		c.Set("Logger", logger)

		rawData, _ := c.GetRawData()

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.String("query_params", c.Request.URL.RawQuery),
			zap.String("body", string(rawData)),
			zap.String("client_ip", c.ClientIP()))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))
		c.Next()
	}
}
