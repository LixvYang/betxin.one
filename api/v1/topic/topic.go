package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/db"
	"github.com/lixvyang/betxin.one/internal/model/redis"
)

type TopicHandler struct {
	db    db.Database
	redis redis.IRedis
}

func NewUserHandler(db db.Database, rds redis.IRedis) *TopicHandler {
	return &TopicHandler{
		db:    db,
		redis: rds,
	}
}

type ITopicHandler interface {
	Create(*gin.Context)
}
