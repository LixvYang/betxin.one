package refund

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
)

type RefundHandler struct {
	storage  *mongo.MongoService
	mixinSrv *mixin_srv.MixinSrv
	cache    *cache.Cache
}

func NewHandler(db *mongo.MongoService, mixinSrv *mixin_srv.MixinSrv, cache *cache.Cache) *RefundHandler {
	return &RefundHandler{
		storage:  db,
		mixinSrv: mixinSrv,
		cache:    cache,
	}
}

type IRefundHandler interface {
}
