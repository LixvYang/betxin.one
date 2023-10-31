package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/db"
	"github.com/lixvyang/betxin.one/internal/model/redis"
)

type UserHandler struct {
	db    db.User
	redis redis.IRedis
}

func NewUserHandler(db db.User, rds redis.IRedis) *UserHandler {
	return &UserHandler{
		db:    db,
		redis: rds,
	}
}

type IUserHandler interface {
	Create(*gin.Context)
}
