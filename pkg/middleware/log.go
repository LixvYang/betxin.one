package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GinLogger(Lg *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		defer func() {
			now := time.Since(start).Milliseconds()

			Lg.Info().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("query", query).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
				Int64("cost(ms)", now).Send()
		}()

		c.Next()
	}
}
