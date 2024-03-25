package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/rs/zerolog"
)

// type Adapter struct {
// 	pool sync.Pool
// }

// func New() *Adapter {
// 	return &Adapter{
// 		pool: sync.Pool{
// 			New: func() interface{} {
// 				return bytes.NewBuffer(make([]byte, 4096))
// 			},
// 		},
// 	}
// }

// var adapter *Adapter

// func init() {
// 	adapter = New()
// }

// 重构读取 body 内容
func GinLogger(Lg *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		defer func() {
			logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

			// buffer := adapter.pool.Get().(*bytes.Buffer)
			// buffer.Reset()
			// defer func() {
			// 	if buffer != nil {
			// 		adapter.pool.Put(buffer)
			// 		buffer = nil
			// 	}
			// }()

			// _, err := io.Copy(buffer, c.Request.Body)
			// if err != nil {
			// 	logger.Error().Err(err).Send()
			// 	return
			// }

			logger.Info().
				// Str("body", string(body)).
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
