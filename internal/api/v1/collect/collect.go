package collect

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
)

type CollectHandler struct {
	storage *mongo.MongoService
}

func NewHandler(db *mongo.MongoService) *CollectHandler {
	return &CollectHandler{
		storage: db,
	}
}

type ICollectHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
}
