package bonuse

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type BonuseModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewBonuseModel(query *query.Query, cache *cache.Cache) BonuseModel {
	return BonuseModel{
		db:    query,
		cache: cache,
	}
}

func (bm *BonuseModel) CreateBonuse(*schema.Bonuse) error {
	return nil
}

func (bm *BonuseModel) GetBonuseByTraceId(string) (*schema.Bonuse, error) {
	return nil, nil
}

func (bm *BonuseModel) ListBonuses() ([]*schema.Bonuse, error) {
	return nil, nil
}

func (bm *BonuseModel) UpdateBonuse(*schema.Bonuse) error {
	return nil
}

func (bm *BonuseModel) DeleteBonuse(string) error {
	return nil
}

func (bm *BonuseModel) GetBonusesByUid(string) (*schema.Bonuse, error) {
	return nil, nil
}
