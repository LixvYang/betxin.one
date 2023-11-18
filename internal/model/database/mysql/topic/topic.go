package topic

import (
	"context"
	"encoding/json"
	"fmt"
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

type HandlerFunc func(ctx context.Context, logger *zerolog.Logger, q *query.Query, sqlTopics ...*sqlmodel.Topic) error

type TopicModel struct {
	db       *query.Query
	cache    *cache.Cache
	handlers map[string][]HandlerFunc
}

func (um *TopicModel) Use(name string, f HandlerFunc) {
	if um.handlers == nil {
		um.handlers = make(map[string][]HandlerFunc)
	}
	um.handlers[name] = append(um.handlers[name], f)
}

func (um *TopicModel) Run(ctx context.Context, logger *zerolog.Logger, name string, sqlTopics ...*sqlmodel.Topic) error {
	var err error
	for _, f := range um.handlers[name] {
		if err = f(ctx, logger, um.db, sqlTopics...); err != nil {
			return err
		}
	}
	return nil
}

func NewTopicModel(q *query.Query, cache *cache.Cache) TopicModel {
	tm := TopicModel{
		db:    q,
		cache: cache,
	}

	tm.Use("AfterFind", tm.AfterFind)
	tm.Use("BeforeUpdate", tm.BeforeUpdate)
	return tm
}

func (um *TopicModel) StopTopic(ctx context.Context, logger *zerolog.Logger, tid int64) error {
	// 删除话题缓存
	um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", tid))
	defer func() {
		time.Sleep(consts.DelayedDeletionInterval)
		um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", tid))
	}()

	sqlTopic, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid), query.Topic.IsDeleted.Zero()).Last()
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in mysql", tid)
		return err
	}
	// err = um.checkUpdate(originTopic)
	// if err != nil {
	// 	return err
	// }
	err = um.Run(ctx, logger, "BeforeUpdate", sqlTopic)
	if err != nil {
		return err
	}

	_, err = um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Update(query.Topic.IsStop, true)
	if err != nil {
		return err
	}

	return nil
}

func (um *TopicModel) CheckTopicExist(ctx context.Context, logger *zerolog.Logger, tid int64) error {
	_, err := um.GetTopicByTid(ctx, logger, tid)
	if err != nil {
		return err
	}
	return nil
}

func (um *TopicModel) CheckTopicStop(ctx context.Context, logger *zerolog.Logger, tid int64) error {
	sqlTopic, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Last()
	if err != nil {
		return err
	}
	if sqlTopic.IsStop {
		return errors.New("topic already stop")
	}

	return nil
}

func (um *TopicModel) GetTopicsByCid(ctx context.Context, logger *zerolog.Logger, cid int64) (topics []*schema.Topic, err error) {
	sqlTopics, err := um.db.Topic.WithContext(ctx).Where(um.db.Topic.Cid.Eq(cid)).Find()
	if err != nil {
		return nil, err
	}

	um.Run(ctx, logger, "AfterFind", sqlTopics...)
	copier.Copy(&topics, &sqlTopics)
	return topics, nil
}

func (um *TopicModel) GetTopicByTid(ctx context.Context, logger *zerolog.Logger, tid int64) (topic *schema.Topic, err error) {
	topic, err = um.getTopicinfoFromCache(ctx, logger, tid)
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in cache", tid)
	} else {
		return topic, err
	}
	sqlTopic, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid), query.Topic.IsDeleted.Zero()).Last()
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in mysql", tid)
		return nil, err
	}

	topic = new(schema.Topic)
	copier.Copy(topic, sqlTopic)
	um.Run(ctx, logger, "AfterFind", sqlTopic)
	go um.encodeTopicInfoToCache(ctx, logger, topic)

	return topic, nil
}

func (um *TopicModel) CreateTopic(ctx context.Context, logger *zerolog.Logger, topic *schema.Topic) error {
	sqlTopic := new(sqlmodel.Topic)
	copier.Copy(sqlTopic, topic)
	return um.db.Topic.WithContext(ctx).Debug().Create(sqlTopic)
}

func (um *TopicModel) DeleteTopic(ctx context.Context, logger *zerolog.Logger, tid int64) (err error) {
	// 延时双删除
	defer func() {
		go func() {
			time.Sleep(consts.DelayedDeletionInterval)
			um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", tid))
		}()
	}()
	// 缓存找
	_, err = um.getTopicinfoFromCache(ctx, logger, tid)
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in cache", tid)
	} else {
		// 删缓存
		um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", tid))
	}

	// 数据库找
	_, err = um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Last()
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in mysql", tid)
		return
	}
	// 数据库删除数据
	_, err = um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(tid)).Update(query.Topic.IsDeleted, true)
	if err != nil {
		logger.Info().Msgf("tid: %d, delete failed in mysql", tid)
		return
	}

	return nil
}

func (um *TopicModel) UpdateTopicInfo(ctx context.Context, logger *zerolog.Logger, topic *schema.Topic) error {
	// 删除缓存
	um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", topic.Tid))

	// 延时双删
	defer func() {
		um.cache.HDel(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", topic.Tid))
	}()

	sqlTopic, err := um.db.Topic.WithContext(ctx).Where(query.Topic.Tid.Eq(topic.Tid), query.Topic.IsDeleted.Zero()).Last()
	if err != nil {
		logger.Info().Msgf("tid: %d, not found in mysql", topic.Tid)
		return err
	}

	err = um.Run(ctx, logger, "BeforeUpdate", sqlTopic)
	if err != nil {
		return err
	}
	// err = um.checkUpdate(originTopic)
	// if err != nil {
	// 	return err
	// }

	// 开启事务
	tx := um.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			panic(err)
		}
	}()

	t := query.Topic
	_, err = t.WithContext(ctx).
		Where(t.Tid.Eq(topic.Tid)).
		UpdateColumns(&sqlmodel.Topic{
			Tid:           topic.Tid,
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// TODO 更新话题价格应该在一个事务中进行
// 并且涉及转账支付业务
func (um *TopicModel) UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
}

func (um *TopicModel) ListTopicByCid(ctx context.Context, logger *zerolog.Logger, cid int64, preId int64, pageSize int64) (topics []*schema.Topic, err error) {
	sqlTopics, err := um.db.WithContext(ctx).Topic.Where(query.Topic.Cid.Eq(cid), query.Topic.Tid.Lte(preId), query.Topic.DeletedAt.Eq(0), query.Topic.IsDeleted.Zero()).Limit(int(pageSize)).Order(query.Topic.Tid.Desc()).Find()
	if err != nil {
		return nil, err
	}
	um.Run(ctx, logger, "AfterFind", sqlTopics...)
	copier.Copy(&topics, &sqlTopics)

	return topics, nil
}

func (um *TopicModel) encodeTopicInfoToCache(ctx context.Context, logger *zerolog.Logger, data *schema.Topic) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error().Msgf("encode node to bytes fail, %+v", data)
		return
	}
	err = um.cache.HSet(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", data.Tid), bytes)
	if err != nil {
		logger.Error().Err(err).Msgf("encode topic to redis fail, %+v", data)
	}
}

func (um *TopicModel) getTopicinfoFromCache(ctx context.Context, logger *zerolog.Logger, tid int64) (*schema.Topic, error) {
	bytes, err := um.cache.HGet(ctx, consts.RdsHashTopicInfoKey, fmt.Sprintf("%d", tid))
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

// func (um *TopicModel) checkUpdate(t *schema.Topic) (err error) {
// 	if t.IsDeleted {
// 		return errors.New("topic already deleted")
// 	}

// 	if t.IsStop || time.Now().After(time.UnixMilli(t.EndTime)) {
// 		return errors.New("topic already stop")
// 	}
// 	decimal.DivisionPrecision = 2
// 	yesCnt, _ := decimal.NewFromString(t.YesCount)

// 	totalCnt, err := decimal.NewFromString(t.TotalCount)
// 	if err != nil {
// 		return err
// 	}
// 	if totalCnt.Equal(decimal.NewFromFloat(0)) {
// 		return nil
// 	}
// 	yesRatio := yesCnt.Div(totalCnt)
// 	t.YesRatio = yesRatio.String()
// 	t.NoRatio = decimal.NewFromInt(100).Sub(yesRatio).String()
// 	return nil
// }
