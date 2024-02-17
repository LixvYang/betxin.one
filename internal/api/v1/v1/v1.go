package v1

import (
	"context"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/api/v1/category"
	"github.com/lixvyang/betxin.one/internal/api/v1/topic"
	"github.com/lixvyang/betxin.one/internal/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	user.IUserHandler
	// collect.ICollectHandler
	topic.ITopicHandler
	category.ICategoryHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig) *BetxinHandler {
	db := database.New(logger, conf)
	// userz := userz.NewSrv(db)
	cache := cache.New(logger, conf.RedisConfig)

	go func() {
		db.CreateCategory(context.Background(), logger, &schema.Category{
			ID:   1,
			Name: "测试",
		})
	}()

	return &BetxinHandler{
		IUserHandler:  user.NewHandler(db, logger),
		ITopicHandler: topic.NewHandler(db, cache),
		// ICollectHandler: collect.NewHandler(db),
		ICategoryHandler: category.NewHandler(db),
	}
}
