package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger is a gin common logging middleware.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger = InjectedLogger(c.Request.Context(), logger)
		start := time.Now()

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)

		reqFields := []zapcore.Field{
			zap.String("proto", c.Request.Proto),
			zap.String("uri", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.String("remote", c.Request.RemoteAddr),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("request", string(body)),
		}

		logger.Debug(
			"request started",
			reqFields...,
		)

		ww := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = ww

		c.Next()

		respFields := []zapcore.Field{
			zap.Int("status", ww.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("response", ww.body.String()),
		}

		logger.Debug(
			"response",
			respFields...,
		)
	}
}
