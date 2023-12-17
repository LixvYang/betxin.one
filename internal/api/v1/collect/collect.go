package collect

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
)

type CollectHandler struct {
	collectStore core.CollectStore
	// cache        *cache.Cache
}

func NewHandler(db database.Database) ICollectHandler {
	collectHandler := &CollectHandler{
		collectStore: db,
		// cache:        cache,
	}

	return collectHandler
}

type ICollectHandler interface {
	Create(*gin.Context)
	Delete(*gin.Context)
}
