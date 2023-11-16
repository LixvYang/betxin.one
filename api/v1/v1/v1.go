package v1

import (
	"github.com/lixvyang/betxin.one/api/v1/topic"
	"github.com/lixvyang/betxin.one/api/v1/user"
	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
}

func NewBetxinHandler(conf *configs.AppConfig) *BetxinHandler {
	db := database.New(conf)

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
	}
}
