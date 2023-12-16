package category

import (
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/category/dao"
	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: store.OutDirPrefix + "category/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.CategoryStore) {}, core.Category{})
		},
	)
}

func New(h *store.Handler) core.CategoryStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Category).(core.CategoryStore)
	if !ok {
		panic("dao.Category is not core.CategoryStore")
	}

	return &storeImpl{
		CategoryStore: v,
	}
}

type storeImpl struct {
	core.CategoryStore
}
