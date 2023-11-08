package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type TopicHandler struct {
	storage database.ITopic
}

func NewHandler(db database.ITopic) ITopicHandler {
	return &TopicHandler{
		storage: db,
	}
}

type ITopicHandler interface {
	Create(*gin.Context)
}
