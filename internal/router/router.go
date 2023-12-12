package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/sd"
	"github.com/lixvyang/betxin.one/api/v1/v1"
	configs "github.com/lixvyang/betxin.one/config"
	_ "github.com/lixvyang/betxin.one/docs"
	"github.com/lixvyang/betxin.one/pkg/middleware"
	"github.com/rs/zerolog"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	server *http.Server
}

func NewService(logger *zerolog.Logger, conf *configs.AppConfig) *Service {
	router := initRouter(logger, conf)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Msgf("listen: %s+v\n", err)
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

func initRouter(logger *zerolog.Logger, conf *configs.AppConfig) *gin.Engine {
	e := gin.New()

	e.Use(
		middleware.Cors(),
		middleware.GinXid(logger),
		middleware.GinLogger(logger),
		middleware.GinRecovery(logger, true),
	)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	h := v1.NewBetxinHandler(logger, conf)
	api := e.Group("/api/v1")
	{
		// 用户
		api.POST("/connect", h.IUserHandler.Connect)

		// 话题查询
		api.GET("/topics/:cid", h.ITopicHandler.ListTopicsByCid)
		api.POST("/topic", h.ITopicHandler.Create)
		api.PUT("/topic/:tid", h.ITopicHandler.UpdateTopicInfo)
		api.DELETE("/topic/:tid", h.ITopicHandler.Delete)
		api.GET("/topic/:tid", h.ITopicHandler.Get)
	}

	// 管理员权限
	{
		api.POST("/category", h.ICategoryHandler.Create)
		api.DELETE("/category/:id", h.ICategoryHandler.Delete)
		api.PUT("/category/:id", h.ICategoryHandler.Update)
		api.GET("/category/:id", h.ICategoryHandler.Get)
		api.GET("/categories", h.ICategoryHandler.List)
	}

	{
		api.GET("/backend/health", sd.HealthCheck)
		api.GET("/backend/disk", sd.DiskCheck)
		api.GET("/backend/cpu", sd.CPUCheck)
		api.GET("/backend/ram", sd.RAMCheck)
	}

	api.Use(middleware.JWTAuthMiddleware())
	{
		api.GET("/user", h.IUserHandler.Get)

	}

	return e
}
