package middleware

import (
	"go-template/internal/version"

	"github.com/gin-gonic/gin"
)

// Version is a gin common middleware to set version in HTTP Header.
func Version() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", version.VERSION)
		c.Header("X-API-Revision", version.REVISION)
		c.Next()
	}
}
