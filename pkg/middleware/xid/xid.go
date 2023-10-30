package xid

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

func GinXid(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		xid := xid.New().String()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("xid", xid)
		})
		c.Header("xid", xid)
		c.Set("logger", logger)
		c.Next()
	}
}
