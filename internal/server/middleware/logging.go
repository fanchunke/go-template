package middleware

import (
	"bytes"
	"fmt"
	"go-template/internal/log"
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
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Ctx(c.Request.Context())
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
		}

		logger.Debug(
			fmt.Sprintf("request: %s", string(body)),
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
		}

		logger.Debug(
			fmt.Sprintf("response: %s", ww.body.String()),
			respFields...,
		)
	}
}
