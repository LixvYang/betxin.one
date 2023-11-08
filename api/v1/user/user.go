package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type UserHandler struct {
	storage database.IUser
	redis   *cache.Cache
}

func NewHandler(db database.IUser) IUserHandler {
	return &UserHandler{
		storage: db,
	}
}

type IUserHandler interface {
	Create(*gin.Context)
}
