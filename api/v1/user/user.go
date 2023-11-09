package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type UserHandler struct {
	storage database.IUser
}

func NewHandler(db database.IUser) IUserHandler {
	return &UserHandler{
		storage: db,
	}
}

type IUserHandler interface {
	Connect(c *gin.Context)
}
