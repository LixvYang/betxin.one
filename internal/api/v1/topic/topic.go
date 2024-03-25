package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
)

type TopicHandler struct {
	storage  *mongo.MongoService
	cache    *cache.Cache
	mixinSrv *mixin_srv.MixinSrv
}

func NewHandler(db *mongo.MongoService, cache *cache.Cache, mixinSrv *mixin_srv.MixinSrv) *TopicHandler {
	topicHandler := &TopicHandler{
		storage:  db,
		cache:    cache,
		mixinSrv: mixinSrv,
	}

	return topicHandler
}

type ITopicHandler interface {
	ListTopics(*gin.Context)
	Create(*gin.Context)
	Delete(*gin.Context)
	Get(*gin.Context)
	ListTopicsByCid(*gin.Context)
	UpdateTopicInfo(*gin.Context)
}
