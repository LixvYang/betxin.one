package v1

import (
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/api/v1/category"
	"github.com/lixvyang/betxin.one/internal/api/v1/collect"
	"github.com/lixvyang/betxin.one/internal/api/v1/topic"
	"github.com/lixvyang/betxin.one/internal/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
	category.ICategoryHandler
	collect.ICollectHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig) *BetxinHandler {
	db := database.New(logger, conf)
	cache := cache.New(logger, conf.RedisConfig)

	return &BetxinHandler{
		IUserHandler:     user.NewHandler(db, logger),
		ITopicHandler:    topic.NewHandler(db, cache),
		ICategoryHandler: category.NewHandler(db),
		ICollectHandler:  collect.NewHandler(db),
	}
}
