package v1

import (
	"github.com/lixvyang/betxin.one/api/v1/topic"
	"github.com/lixvyang/betxin.one/api/v1/user"
	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
}

func NewBetxinHandler(conf *configs.AppConfig) *BetxinHandler {
	db := database.New(conf)

	return &BetxinHandler{
		user.NewHandler(db),
		topic.NewHandler(db),
	}
}
