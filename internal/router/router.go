package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/config"
	_ "github.com/lixvyang/betxin.one/docs"
	"github.com/lixvyang/betxin.one/internal/api/sd"
	v1 "github.com/lixvyang/betxin.one/internal/api/v1"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/pkg/middleware"
	"github.com/rs/zerolog"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	server *http.Server
}

func NewService(logger *zerolog.Logger, conf *config.AppConfig) *Service {
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

func initRouter(logger *zerolog.Logger, conf *config.AppConfig) *gin.Engine {
	e := gin.New()

	e.Use(
		middleware.Cors(),
		middleware.GinXid(logger),
		middleware.GinLogger(logger),
		middleware.GinRecovery(logger, true),
	)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 初始化 db 缓存
	db := database.New(logger, conf)
	cache := cache.New(logger, conf.RedisConfig)

	h := v1.NewBetxinHandler(logger, conf, db, cache)
	api := e.Group("/api/v1")
	// api.Use(middleware.JWTAuthMiddleware())
	{
		// 用户
		api.POST("/connect", h.IUserHandler.Connect)

		// 话题相关
		api.GET("/topics", h.ITopicHandler.ListTopics) // 根据query获取话题列表
		api.GET("/topic/:tid", h.ITopicHandler.GetTopicByTid)

		api.POST("/topicpurchase/simulate", h.ITopicPurchaseHandler.Simulate)                           // 提供一个根据传入代币数量 模拟计算最终结果的接口
		api.POST("/topicpurchase/action", h.ITopicPurchaseHandler.CreatePurchaseAction)                 // 话题购买行为表
		api.GET("/topicpurchasehistory/:request_id", h.ITopicPurchaseHandler.QueryTopicPurchaseHistory) // 第二个接口供前端不断判断 是否已购买

		// 话题退款
		api.POST("/refund/topic/", h.IRefundHandler.CreateTopicRefund)

		// 即将结束的话题

		// 话题购买 其实就是传入uuid 然后根据uuid 查询 话题购买行为表 然后根据行为表 计算退款金额

		//收藏相关
		// 话题收藏
		// 取消收藏
		// 用户的话题收藏列表

		// 收藏相关
		demoAPI := api.Use(middleware.DemoAuthNotMiddleware())
		{
			demoAPI.POST("/collect", h.ICollectHandler.Create)
			demoAPI.DELETE("/collect", h.ICollectHandler.Delete)
		}
		// 用户信息
		api.GET("/user", h.IUserHandler.GetUserInfo)
		// 用户信息修改
	}

	admin := e.Group("/admin")
	// admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.POST("/topic", h.ITopicHandler.Create)
		admin.PUT("/topic/:tid", h.ITopicHandler.UpdateTopicInfo)
		admin.DELETE("/topic/:tid", h.ITopicHandler.Delete)

		// 处理话题停止
		admin.POST("/topic/action/stop", h.ITopicHandler.StopTopic)

		admin.POST("/category", h.ICategoryHandler.Create)
		admin.DELETE("/category/:id", h.ICategoryHandler.Delete)
		admin.PUT("/category/:id", h.ICategoryHandler.Update)
		admin.GET("/category/:id", h.ICategoryHandler.Get)
		admin.GET("/categories", h.ICategoryHandler.List)

		admin.GET("/backend/health", sd.HealthCheck)
		admin.GET("/backend/disk", sd.DiskCheck)
		admin.GET("/backend/cpu", sd.CPUCheck)
		admin.GET("/backend/ram", sd.RAMCheck)
	}

	return e
}
