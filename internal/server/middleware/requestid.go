package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type correlationIDType int

const (
	requestIDHeader                   = "X-Request-ID"
	requestIDKey    correlationIDType = iota
	sessionIDKey
)

// RequestID is a middleware that injects a request ID into the context of each
// request. context is `context.Context`, not `gin.Context`
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}

// WithRequestID injects requestID with the given ctx.
func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

// WithSessionID injects sessionID with the given ctx.
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

// GetRequestID returns requestID injected into the context of request.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// InjectedLogger return a *zap.Logger with requestID and sessionID added.
func InjectedLogger(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if ctx != nil {
		requestID := GetRequestID(ctx)
		if requestID != "" {
			logger = logger.With(zap.String("requestID", requestID))
		}
		if sessionID, ok := ctx.Value(sessionIDKey).(string); ok {
			logger = logger.With(zap.String("sessionID", sessionID))
		}
	}
	return logger
}
