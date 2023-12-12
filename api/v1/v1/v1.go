package v1

import (
	"github.com/lixvyang/betxin.one/api/v1/bonuse"
	"github.com/lixvyang/betxin.one/api/v1/category"
	"github.com/lixvyang/betxin.one/api/v1/topic"
	"github.com/lixvyang/betxin.one/api/v1/user"
	configs "github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
	category.ICategoryHandler
	bonuse.IBonuseHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *configs.AppConfig) *BetxinHandler {
	db := database.New(logger, conf)

	categorys, err := db.ListCategories()
	if err != nil {
		panic(err)
	}
	categoryMap := make(map[int64]*schema.Category)
	for _, category := range categorys {
		categoryMap[category.ID] = category
	}

	return &BetxinHandler{
		user.NewHandler(db),
		topic.NewHandler(db, categoryMap),
		category.NewHandler(db),
		bonuse.NewHandler(db),
	}
}
