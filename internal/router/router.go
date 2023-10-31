package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/lixvyang/betxin.one/pkg/middleware"
	"github.com/rs/zerolog"

	v1 "github.com/lixvyang/betxin.one/api/v1"
)

type Service struct {
	server *http.Server
}

func NewService() *Service {
	router := InitRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Conf.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Lg.Fatal().Msgf("listen: %s+v\n", err)
		}
	}()

	return &Service{srv}
}

func (srv *Service) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.server.Shutdown(ctx)
}

func (srv *Service) ListenAndServe() error {
	return srv.server.ListenAndServe()
}

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.GinXid(&logger.Lg,
			middleware.WithLoggerKey(consts.BetxinLoggerKey),
			middleware.WithXid(consts.BetxinXid)),
		middleware.GinLogger(&logger.Lg),
		middleware.GinRecovery(&logger.Lg, true),
	)

	btxHandler := v1.NewBetxinHandler()
	g := r.Group("/api/v1")
	{
		g.POST("/user", btxHandler.ITopicHandler.Create)
	}

	r.GET("/hello", func(c *gin.Context) {
		xl := c.MustGet(consts.BetxinLoggerKey).(*zerolog.Logger)
		xl.Info().Msg("Hello world")

		c.JSON(200, gin.H{
			"Hello": "World",
		})
	})

	r.GET("/world", func(c *gin.Context) {
		xl := c.MustGet(consts.BetxinLoggerKey).(zerolog.Logger)
		xl.Info().Msg("Hello world")

		panic("123")
	})

	return r
}
