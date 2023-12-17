package collect

import (
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/collect/dao"
	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: store.OutDirPrefix + "collect/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.CollectStore) {}, core.Collect{})
		},
	)
}

func New(h *store.Handler) core.CollectStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Collect).(core.CollectStore)
	if !ok {
		panic("dao.collect is not core.UserStore")
	}

	return &storeImpl{
		CollectStore: v,
	}
}

type storeImpl struct {
	core.CollectStore
}
