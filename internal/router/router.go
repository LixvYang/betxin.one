package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/sd"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/lixvyang/betxin.one/pkg/middleware"
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
	e := gin.New()

	e.Use(
		middleware.GinXid(&logger.Lg, middleware.WithLoggerKey(consts.LoggerKey)),
		middleware.GinLogger(&logger.Lg),
		middleware.GinRecovery(&logger.Lg, true),
	)

	btxHandler := handler.NewBetxinHandler()
	betxinAPI := e.Group("/api/v1")
	{
		betxinAPI.POST("/user", btxHandler.IUserHandler.Create)
	}

	{
		betxinAPI.GET("/backend/health", sd.HealthCheck)
		betxinAPI.GET("/backend/disk", sd.DiskCheck)
		betxinAPI.GET("/backend/cpu", sd.CPUCheck)
		betxinAPI.GET("/backend/ram", sd.RAMCheck)
	}

	return e
}
