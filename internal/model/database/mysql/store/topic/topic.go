package topic

import (
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/topic/dao"
	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: store.OutDirPrefix + "topic/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.TopicStore) {}, core.Topic{})
		},
	)
}

func New(h *store.Handler) core.TopicStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Topic).(core.TopicStore)
	if !ok {
		panic("dao.Topic is not core.UserStore")
	}

	return &storeImpl{
		TopicStore: v,
	}
}

type storeImpl struct {
	core.TopicStore
}
