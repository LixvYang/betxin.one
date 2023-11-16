package snapshot

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type SnapshotModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewMessageModel(query *query.Query, cache *cache.Cache) SnapshotModel {
	return SnapshotModel{
		db:    query,
		cache: cache,
	}
}

func (sm *SnapshotModel) CreateSnapshot(*schema.Snapshot) error {
	return nil
}

func (sm *SnapshotModel) ListMessage(uid string) ([]*schema.Snapshot, error) {
	return nil, nil
}

func (sm *SnapshotModel) DeleteSnapshot(trace_id string) error {
	return nil
}

func (sm *SnapshotModel) GetSnapshot(trace_id string) (*schema.Snapshot, error) {
	return nil, nil
}
