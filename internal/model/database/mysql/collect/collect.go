package collect

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

/*
CheckCollect(uid, tid string) (*schema.Collect, error)
	CreateCollect(uid, tid string) error
	ListCollects() ([]*schema.Collect, error)
	GetCollectByUserId(uid string) ([]*schema.Collect, error)
	UpdateCollect(uid, tid string, status bool) (*schema.Collect, error)

*/

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

func (cm *CollectModel) CheckCollect(uid, tid string) (*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) CreateCollect(uid, tid string) error {
	return nil
}

func (cm *CollectModel) ListCollects() ([]*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) GetCollectByUid(uid string) ([]*schema.Collect, error) {
	return nil, nil
}

func (cm *CollectModel) UpdateCollect(uid, tid string, status bool) (*schema.Collect, error) {
	return nil, nil
}
