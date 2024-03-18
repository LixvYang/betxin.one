package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
)

type TopicHandler struct {
	topicSrv    database.ITopic
	collectSrv  database.ICollect
	categorySrv database.ICategory
	cache       *cache.Cache
}

func NewHandler(db database.Database, cache *cache.Cache) ITopicHandler {
	topicHandler := &TopicHandler{
		topicSrv:    db,
		collectSrv:  db,
		categorySrv: db,
		cache:       cache,
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

// type CategorySrv struct {
// 	categorySrv database.ICategory
// }

// func NewCategorySrv(categorySrv database.Database) *CategorySrv {
// 	return &CategorySrv{
// 		categorySrv: categorySrv,
// 	}
// }

// func (c *CategorySrv) Get(id int64) (*schema.Category, error) {
// 	category, err := c.categorySrv.GetCategoryById(context.Background(), &log.Logger, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return category, nil
// }
