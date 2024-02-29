package v1

import (
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/api/v1/category"
	"github.com/lixvyang/betxin.one/internal/api/v1/collect"
	"github.com/lixvyang/betxin.one/internal/api/v1/topic"
	"github.com/lixvyang/betxin.one/internal/api/v1/user"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
	"github.com/rs/zerolog"
)

type BetxinHandler struct {
	mixinSrv *mixin_srv.MixinCli

	user.IUserHandler
	topic.ITopicHandler
	category.ICategoryHandler
	collect.ICollectHandler
}

func NewBetxinHandler(logger *zerolog.Logger, conf *config.AppConfig, db database.Database, cache *cache.Cache) *BetxinHandler {
	mixinSrv := mixin_srv.New(conf.MixinConfig, db)
	go func() {
		// mixinSrv.ArrgegateUtxos(context.Background())
	}()

	return &BetxinHandler{
		mixinSrv: mixinSrv,

		IUserHandler:     user.NewHandler(db, logger),
		ITopicHandler:    topic.NewHandler(db, cache),
		ICategoryHandler: category.NewHandler(db),
		ICollectHandler:  collect.NewHandler(db),
	}
}
