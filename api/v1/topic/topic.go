package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/redis"
)

type TopicHandler struct {
	db    database.ITopic
	redis *redis.RedisClient
}

func NewUserHandler(db database.ITopic, rds *redis.RedisClient) ITopicHandler {
	return &TopicHandler{
		db:    db,
		redis: rds,
	}
}

type ITopicHandler interface {
	Create(*gin.Context)
}
