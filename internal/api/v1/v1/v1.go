package v1

import (
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/api/v1/topic"
	"github.com/lixvyang/betxin.one/internal/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	user.IUserHandler
	// collect.ICollectHandler
	topic.ITopicHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig) *BetxinHandler {
	db := database.New(logger, conf)
	// userz := userz.NewSrv(db)
	// categoryMap := make(map[int64]*core.Category)
	// cache := cache.New(logger, conf.RedisConfig)

	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	categorys, err := db.ListCategories(context.Background())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for _, category := range categorys {
	// 		categoryMap[category.ID] = category
	// 	}
	// }()

	return &BetxinHandler{
		IUserHandler: user.NewHandler(db, logger),
		// ITopicHandler:   topic.NewHandler(db, categoryMap, cache),
		// ICollectHandler: collect.NewHandler(db),
	}
}
