package v1

import (
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/api/v1/category"
	"github.com/lixvyang/betxin.one/internal/api/v1/collect"
	"github.com/lixvyang/betxin.one/internal/api/v1/refund"
	"github.com/lixvyang/betxin.one/internal/api/v1/topic"
	"github.com/lixvyang/betxin.one/internal/api/v1/topicpurchase"
	"github.com/lixvyang/betxin.one/internal/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	mixinSrv *mixin_srv.MixinSrv

	IUserHandler          *user.UserHandler
	ITopicHandler         *topic.TopicHandler
	ICategoryHandler      *category.CategoryHandler
	ICollectHandler       *collect.CollectHandler
	ITopicPurchaseHandler *topicpurchase.TopicPurchaseHandler
	IRefundHandler        *refund.RefundHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig, mongoSrv *mongo.MongoService, cache *cache.Cache) *BetxinHandler {
	mixinSrv := mixin_srv.New(conf)
	go func() {
		// mixinSrv.ArrgegateUtxos(context.Background())
	}()

	return &BetxinHandler{
		mixinSrv: mixinSrv,

		IUserHandler:     user.NewHandler(mongoSrv, logger),
		ITopicHandler:    topic.NewHandler(mongoSrv, cache, mixinSrv),
		ICategoryHandler: category.NewHandler(mongoSrv),
		ICollectHandler:  collect.NewHandler(mongoSrv),
		IRefundHandler:   refund.NewHandler(mongoSrv, mixinSrv, cache),
	}
}
