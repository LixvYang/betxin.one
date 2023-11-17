package collect

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type CollectModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewCollectModel(query *query.Query, cache *cache.Cache) CollectModel {
	return CollectModel{
		db:    query,
		cache: cache,
	}
}

func (cm *CollectModel) CheckCollect(uid string, tid int64) (*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) CreateCollect(uid string, tid int64) error {
	return nil
}

func (cm *CollectModel) ListCollects() ([]*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) GetCollectByUid(uid string) ([]*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) UpdateCollect(uid string, tid int64, status bool) (*schema.Collect, error) {
	return nil, nil
}
