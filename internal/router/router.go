package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/sd"
	"github.com/lixvyang/betxin.one/api/v1/v1"
	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/lixvyang/betxin.one/pkg/middleware"
)

type Service struct {
	server *http.Server
}

func NewService(conf *configs.AppConfig) *Service {
	router := initRouter(conf)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
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

func initRouter(conf *configs.AppConfig) *gin.Engine {
	e := gin.New()

	e.Use(
		middleware.GinXid(&logger.Lg),
		middleware.GinLogger(&logger.Lg),
		middleware.GinRecovery(&logger.Lg, true),
	)

	h := v1.NewBetxinHandler(conf)
	api := e.Group("/api/v1")
	{
		// 用户
		api.POST("/connect", h.IUserHandler.Connect)
	
		// 话题查询
		api.POST("/topic/:cid", h.ITopicHandler.ListTopicsByCid)

	}

	{
		api.GET("/backend/health", sd.HealthCheck)
		api.GET("/backend/disk", sd.DiskCheck)
		api.GET("/backend/cpu", sd.CPUCheck)
		api.GET("/backend/ram", sd.RAMCheck)
	}

	return e
}
