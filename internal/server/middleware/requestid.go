package middleware

import (
	"go-template/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

// RequestID is a middleware that injects a request ID into the context of each
// request. context is `context.Context`, not `gin.Context`
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		logger := zap.L().With(zap.String(requestIDHeader, requestID))
		ctx := log.NewContext(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}
