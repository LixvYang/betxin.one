package topic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/pkg/errors"
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

func (um *TopicModel) StopTopic(ctx context.Context, logger *zerolog.Logger, tid string) error {
	// 删除话题缓存
	um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, tid)
	defer func() {
		go func() {
			time.Sleep(time.Second / 2)
			um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, tid)
		}()
	}()

	_, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Update(query.Topic.IsStop, true)
	if err != nil {
		return err
	}

	return nil
}

func (um *TopicModel) CheckTopicExist(ctx context.Context, logger *zerolog.Logger, tid string) error {
	_, err := um.GetTopicByTid(ctx, logger, tid)
	if err != nil {
		return err
	}
	return nil
}

func (um *TopicModel) CheckTopicStop(ctx context.Context, logger *zerolog.Logger, tid string) error {
	topic, err := um.GetTopicByTid(ctx, logger, tid)
	if err != nil {
		return err
	}
	if topic.IsStop {
		return errors.New("topic already stop")
	}

	return nil
}

func (um *TopicModel) GetTopicsByCid(ctx context.Context, logger *zerolog.Logger, cid int64) (topics []*schema.Topic, count int, err error) {
	sqlTopics, err := um.db.Topic.WithContext(ctx).Where(um.db.Topic.Cid.Eq(cid)).Find()
	if err != nil {
		return nil, 0, err
	}
	copier.Copy(topics, sqlTopics)

	return topics, len(topics), nil
}

func (um *TopicModel) GetTopicByTid(ctx context.Context, logger *zerolog.Logger, tid string) (topic *schema.Topic, err error) {
	topic = new(schema.Topic)
	topic, err = um.getTopicinfoFromCache(ctx, logger, tid)
	if err != nil {
		logger.Info().Msgf("tid: %s, not found in cache", tid)
	} else {
		return topic, err
	}
	sqlTopic, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Last()
	if err != nil {
		logger.Info().Msgf("tid: %s, not found in mysql", tid)
		return nil, err
	}
	copier.Copy(sqlTopic, topic)

	go um.encodeTopicInfoToCache(ctx, logger, topic)

	return topic, nil
}

func (um *TopicModel) CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
}

func (um *TopicModel) DeleteTopic(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
}

func (um *TopicModel) UpdateTopicInfo(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
}

func (um *TopicModel) UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
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

func (um *TopicModel) getTopicinfoFromCache(ctx context.Context, logger *zerolog.Logger, tid string) (*schema.Topic, error) {
	bytes, err := um.cache.HGet(ctx, consts.RdsHashTopicInfoKey, tid)
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
