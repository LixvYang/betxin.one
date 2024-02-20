package collect

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type CollectHandler struct {
	storage  database.ICollect
	topicSrv database.ITopic
}

func NewHandler(db database.Database) *CollectHandler {
	return &CollectHandler{
		storage:  db,
		topicSrv: db,
	}
}

type ICollectHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
}
