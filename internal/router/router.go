package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/lixvyang/betxin.one/pkg/middleware/log"
	"github.com/lixvyang/betxin.one/pkg/middleware/recovery"
	"github.com/lixvyang/betxin.one/pkg/middleware/xid"
	"github.com/rs/zerolog"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(
		xid.GinXid(&logger.Lg),
		log.GinLogger(&logger.Lg),
		recovery.GinRecovery(&logger.Lg, true),
	)

	r.GET("/hello", func(c *gin.Context) {
		xl := c.MustGet("logger").(*zerolog.Logger)
		xl.Info().Msg("Hello world")

		c.JSON(200, gin.H{
			"Hello": "World",
		})
	})

	r.GET("/world", func(c *gin.Context) {
		xl := c.MustGet("logger").(zerolog.Logger)
		xl.Info().Msg("Hello world")

		panic("123")
	})

	return r
}
