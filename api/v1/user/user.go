package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
)

type UserHandler struct {
	core.UserStore
	core.UserService
}

func NewHandler(userStore core.UserStore, userService core.UserService) IUserHandler {
	return &UserHandler{
		userStore,
		userService,
	}
}

type IUserHandler interface {
	Connect(c *gin.Context)
	Get(c *gin.Context)
}
