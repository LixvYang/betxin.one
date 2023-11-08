package topic

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type TopicModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewTopicModel(query *query.Query, cache *cache.Cache) TopicModel {
	return TopicModel{
		db:    query,
		cache: cache,
	}
}



func (um *TopicModel) encodeTopicInfoToCache(ctx context.Context, logger *zerolog.Logger, data *schema.Topic) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error().Msgf("encode node to bytes fail, %+v", data)
		return
	}
	err = um.cache.HSet(ctx, consts.RdsHashTopicInfoKey, data.Tid, bytes)
	if err != nil {
		logger.Error().Msgf("encode topic to redis fail, %+v", data)
	}
}

func (um *TopicModel) getTopicinfoFromCache(ctx context.Context, logger *zerolog.Logger, uid string) (*schema.Topic, error) {
	bytes, err := um.cache.HGet(ctx, consts.RdsHashTopicInfoKey, uid)
	// 找到了数据
	if err == nil && bytes != nil {
		var data schema.Topic
		json.Unmarshal(bytes, &data)
		return &data, nil
	}
	// 没有找到数据
	if err == cache.Nil {
		return nil, errors.New(consts.CacheNotFound)
	}
	// redis错误
	return nil, err
}
