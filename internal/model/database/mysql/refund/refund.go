package refund

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type RefundModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewMessageModel(query *query.Query, cache *cache.Cache) RefundModel {
	return RefundModel{
		db:    query,
		cache: cache,
	}
}

func (rm *RefundModel) CreateRefund(*schema.Refund) error {
	return nil
}

func (rm *RefundModel) GetRefundByTraceId(string) (*schema.Refund, error) {
	return nil, nil
}

func (rm *RefundModel) ListRefunds() ([]*schema.Refund, error) {
	return nil, nil
}

func (rm *RefundModel) UpdateRefund() {
	return
}

func (rm *RefundModel) DeleteRefund(trace_id string) error {
	return nil
}

func (rm *RefundModel) GetRefundsByUid(uid string) ([]*schema.Refund, error) {
	return nil, nil
}
