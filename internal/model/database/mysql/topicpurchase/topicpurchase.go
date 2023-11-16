package topicpurchase

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type TopicPurchaseModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewTopicPurchaseModel(query *query.Query, cache *cache.Cache) TopicPurchaseModel {
	return TopicPurchaseModel{
		db:    query,
		cache: cache,
	}
}

func (tpm *TopicPurchaseModel) CheckTopicPurchase(uid, tid string) error {
	return nil
}

func (tpm *TopicPurchaseModel) GetTopicPurchase(uid, tid string) (*schema.TopicPurchase, error) {
	return nil, nil
}

func (tpm *TopicPurchaseModel) CreateTopicPurchase(*schema.TopicPurchase) error {
	return nil
}
