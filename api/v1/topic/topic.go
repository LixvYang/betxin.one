package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type TopicHandler struct {
	storage             database.ITopic
	topicCollectHandler database.ICollect
	categoryMap         map[int64]*schema.Category
}

func NewHandler(db database.Database, categoryMap map[int64]*schema.Category) ITopicHandler {
	topicHandler := &TopicHandler{
		storage:             db,
		topicCollectHandler: db,
		categoryMap:         categoryMap,
	}

	return topicHandler
}

type ITopicHandler interface {
	Create(*gin.Context)
	ListTopicsByCid(*gin.Context)
}
