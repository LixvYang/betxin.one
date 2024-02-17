package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/rs/zerolog"
)

func GinLogger(Lg *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		defer func() {
			logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

			logger.Info().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("query", query).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
				Int64("cost(ms)", time.Since(start).Milliseconds()).Send()
		}()

		c.Next()
	}
}
