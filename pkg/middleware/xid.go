package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

const (
	LoggerKey = "logger"
	Xid       = "xid"
)

type OptionFunc func(*XidOptions)

type XidOptions struct {
	loggerKey string
	xid       string
}

func WithLoggerKey(loggerKey string) OptionFunc {
	return func(xo *XidOptions) {
		xo.loggerKey = loggerKey
	}
}

func WithXid(xid string) OptionFunc {
	return func(xo *XidOptions) {
		xo.xid = xid
	}
}

func GinXid(logger *zerolog.Logger, ofs ...OptionFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		xidOptions := &XidOptions{LoggerKey, Xid}
		for _, of := range ofs {
			of(xidOptions)
		}

		xid := xid.New().String()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(xidOptions.xid, xid)
		})

		c.Header(xidOptions.xid, xid)
		c.Set(xidOptions.loggerKey, logger)
		c.Next()
	}
}
