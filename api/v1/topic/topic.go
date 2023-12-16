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
	return nil
	// topicHandler := &TopicHandler{
	// 	storage:             db,
	// 	topicCollectHandler: db,
	// 	categoryMap:         categoryMap,
	// }

	// return topicHandler
}

type ITopicHandler interface {
	Create(*gin.Context)
	Delete(*gin.Context)
	Get(*gin.Context)
	ListTopicsByCid(*gin.Context)
	UpdateTopicInfo(*gin.Context)
}
