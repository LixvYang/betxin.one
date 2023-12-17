package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
)

type TopicHandler struct {
	topicStore   core.TopicStore
	collectStore core.CollectStore
	categoryMap  map[int64]*core.Category
	cache        *cache.Cache
}

func NewHandler(db database.Database, categoryMap map[int64]*core.Category, cache *cache.Cache) ITopicHandler {
	topicHandler := &TopicHandler{
		topicStore:   db,
		collectStore: db,
		categoryMap:  categoryMap,
		cache:        cache,
	}

	return topicHandler
}

type ITopicHandler interface {
	Create(*gin.Context)
	Delete(*gin.Context)
	Get(*gin.Context)
	ListTopicsByCid(*gin.Context)
	UpdateTopicInfo(*gin.Context)
}

