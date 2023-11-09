package topic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
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

func (um *TopicModel) CreateTopic(ctx context.Context, logger *zerolog.Logger, topic *schema.Topic) error {
	var sqlTopic sqlmodel.Topic
	copier.Copy(&sqlTopic, topic)
	return um.db.Topic.WithContext(ctx).Create(&sqlTopic)
}

func (um *TopicModel) DeleteTopic(ctx context.Context, logger *zerolog.Logger, tid string) (err error) {
	// 延时双删除
	defer func() {
		go func() {
			time.Sleep(time.Second * 3)
			um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, tid)
		}()
	}()
	// 缓存找
	_, err = um.getTopicinfoFromCache(ctx, logger, tid)
	if err != nil {
		logger.Info().Msgf("tid: %s, not found in cache", tid)
	} else {
		// 删缓存
		um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, tid)
	}

	// 数据库找
	_, err = um.db.User.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Last()
	if err != nil {
		logger.Info().Msgf("tid: %s, not found in mysql", tid)
		return
	}
	// 数据库删除数据
	_, err = um.db.User.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Delete()
	if err != nil {
		logger.Info().Msgf("tid: %s, delete failed in mysql", tid)
		return
	}

	return nil
}

func (um *TopicModel) UpdateTopicInfo(ctx context.Context, logger *zerolog.Logger, topic *schema.Topic) error {
	// 删除缓存
	um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, topic.Tid)

	// 延时双删
	defer func() {
		go func() {
			time.Sleep(time.Second * 3)
			um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, topic.Tid)
		}()
	}()

	// 更新数据
	um.db.Transaction(func(tx *query.Query) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Info().Any("recover", r).Send()
			}
		}()

		_, err := query.Topic.WithContext(ctx).
			Where(query.Topic.Tid.Eq(topic.Tid)).
			Updates(sqlmodel.Topic{
				Cid:           topic.Cid,
				Title:         topic.Title,
				Intro:         topic.Intro,
				Content:       topic.Content,
				YesRatio:      topic.YesCount,
				NoRatio:       topic.NoCount,
				YesCount:      topic.YesCount,
				NoCount:       topic.NoCount,
				TotalCount:    topic.TotalCount,
				CollectCount:  topic.CollectCount,
				ReadCount:     topic.ReadCount,
				ImgURL:        topic.ImgURL,
				IsStop:        topic.IsStop,
				RefundEndTime: topic.RefundEndTime,
				EndTime:       topic.EndTime,
			})
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}

// TODO 更新话题价格应该在一个事务中进行
// 并且涉及转账支付业务
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
