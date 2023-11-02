package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/redis"
)

type UserHandler struct {
	db    database.IUser
	redis *redis.RedisClient
}

func NewUserHandler(db database.IUser, rds *redis.RedisClient) IUserHandler {
	return &UserHandler{
		db:    db,
		redis: rds,
	}
}

type IUserHandler interface {
	Create(*gin.Context)
}
