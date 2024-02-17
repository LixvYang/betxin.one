package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

func GinXid(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		xid := xid.New().String()

		log := *logger
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(consts.DefaultXid, xid)
		})

		c.Header(consts.DefaultXid, xid)
		c.Set(consts.DefaultLoggerKey, log)
		c.Set(consts.DefaultXid, xid)

		c.Next()
	}
}
