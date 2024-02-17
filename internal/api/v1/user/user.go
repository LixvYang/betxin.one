package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userSrv database.IUser
}

func NewHandler(database database.Database, logger *zerolog.Logger) IUserHandler {
	return &UserHandler{
		database,
	}
}

type IUserHandler interface {
	Connect(c *gin.Context)
	Get(c *gin.Context)
}
