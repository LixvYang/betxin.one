package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userSrv *mongo.MongoService
}

func NewHandler(database *mongo.MongoService, logger *zerolog.Logger) *UserHandler {
	return &UserHandler{
		database,
	}
}

type IUserHandler interface {
	Connect(c *gin.Context)
	Get(c *gin.Context)
}
