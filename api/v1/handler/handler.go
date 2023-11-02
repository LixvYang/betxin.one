package handler

import (
	"github.com/lixvyang/betxin.one/api/v1/topic"
	"github.com/lixvyang/betxin.one/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/redis"
)

type BetxinHandler struct {
	user.IUserHandler
	topic.ITopicHandler
}

func NewBetxinHandler() *BetxinHandler {
	db := database.NewDatabse()
	rds := redis.NewRedisClient()

	return &BetxinHandler{
		user.NewUserHandler(db, rds),
		topic.NewUserHandler(db, rds),
	}
}
