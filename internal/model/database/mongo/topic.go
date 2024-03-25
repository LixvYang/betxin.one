package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/qiniu/qmgo/options"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var zone, _ = time.LoadLocation("Asia/Shanghai")

func changeTopicTime(topics ...*schema.Topic) {
	lo.ForEach(topics, func(item *schema.Topic, index int) {
		item.EndTime = item.EndTime.In(zone)
		item.RefundEndTime = item.RefundEndTime.In(zone)
		item.CreatedAt = item.CreatedAt.In(zone)
		item.UpdatedAt = item.UpdatedAt.In(zone)
		item.DeletedAt = item.DeletedAt.In(zone)
	})
}

// 寻找话题时不查找 is_stop 的话题
func filterNoStopTopics(filter bson.M) {
	filter["is_stop"] = bson.M{"$eq": false}
}

// 寻找话题时不查找 deleted_at 的话题
func filterDeletedAtTopics(filter bson.M) {
	filter["deleted_at"] = bson.M{"$eq": time.Time{}}
}

// 列出所有未停止 未被删除的话题
// 按结束时间倒序排列
func (s MongoService) ListAllTopics(ctx context.Context) ([]*schema.Topic, int64, error) {
	var topics []*schema.Topic
	var total int64
	var err error

	filter := bson.M{"is_stop": false}
	filterDeletedAtTopics(filter)

	find := s.topicColl.Find(ctx, filter)
	total, err = find.Count()
	if err != nil {
		return nil, total, err
	}

	err = find.Sort("-end_time").All(&topics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, total, ErrNoSuchItem
		}
		return nil, total, err
	}
	changeTopicTime(topics...)
	return topics, total, nil
}

// cid 为 0 表示查找所有的话题
func (s *MongoService) ListTopics(ctx context.Context, cid int64, limit, offset int64) ([]*schema.Topic, int64, error) {
	var topics []*schema.Topic
	var total int64
	var err error

	filter := bson.M{"is_stop": false}
	filterDeletedAtTopics(filter)
	if cid != 0 {
		filter["cid"] = cid
	}

	find := s.topicColl.Find(ctx, filter)
	total, err = find.Count()
	if err != nil {
		return nil, total, err
	}

	err = find.Limit(limit).Skip(offset).Sort("-created_at").All(&topics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, total, ErrNoSuchItem
		}
		return nil, total, err
	}
	changeTopicTime(topics...)
	return topics, total, nil
}

func (s *MongoService) ListTopicByCid(ctx context.Context, cid int64, createdAt time.Time, pageSize int64) ([]*schema.Topic, int64, error) {
	var topics []*schema.Topic
	var total int64
	var err error

	filter := bson.M{"created_at": bson.M{"$lte": createdAt}, "cid": cid, "is_stop": false}
	filterDeletedAtTopics(filter)

	find := s.topicColl.Find(ctx, filter)
	total, err = find.Count()
	if err != nil {
		return nil, total, err
	}

	if total == 0 {
		return nil, total, ErrNoSuchItem
	}
	err = find.Limit(pageSize).Sort("-created_at").All(&topics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, total, ErrNoSuchItem
		}
		return nil, total, err
	}
	changeTopicTime(topics...)

	return topics, total, nil
}

func (s *MongoService) StopTopic(ctx context.Context, tid string) error {
	err := s.topicColl.UpdateOne(ctx, bson.M{"tid": tid, "is_stop": false}, bson.M{"$set": bson.M{"is_stop": true}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}

func (s *MongoService) GetTopicsByCid(ctx context.Context, cid int64) ([]*schema.Topic, error) {
	var topics []*schema.Topic
	filter := bson.M{"cid": cid, "is_stop": false}
	filterDeletedAtTopics(filter)
	err := s.topicColl.Find(ctx, filter).All(&topics)
	if err != nil {
		return nil, err
	}
	changeTopicTime(topics...)

	return topics, nil
}

func (s *MongoService) GetTopicByTid(ctx context.Context, tid string) (*schema.Topic, error) {
	topic := &schema.Topic{}
	filter := bson.M{"tid": tid, "is_stop": false}
	filterDeletedAtTopics(filter)

	err := s.topicColl.Find(ctx, filter, options.FindOptions{
		QueryHook: topic,
	}).One(&topic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	changeTopicTime(topic)
	return topic, nil
}

func (s *MongoService) CreateTopic(ctx context.Context, topic *schema.Topic) error {
	fmt.Println("创建话题tid", topic.Tid)

	_, err := s.topicColl.InsertOne(ctx, topic)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) DeleteTopic(ctx context.Context, tid string) error {
	err := s.topicColl.UpdateOne(ctx, bson.M{"tid": tid}, bson.M{"$set": bson.M{"deleted_at": time.Now().Local()}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}

func (s *MongoService) UpdateTopic(ctx context.Context, tid string, topic *schema.Topic) error {
	filter := bson.M{"tid": tid}
	topic.UpdatedAt = time.Now()

	filterDeletedAtTopics(filter)
	err := s.topicColl.UpdateOne(ctx, filter, bson.M{"$set": topic})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}

func (s *MongoService) GetTopicsByTids(ctx context.Context, tids []string) ([]*schema.Topic, error) {
	var topics []*schema.Topic
	filiter := bson.M{"tid": bson.M{"$in": tids}}
	filterDeletedAtTopics(filiter)
	err := s.topicColl.Find(ctx, filiter).All(&topics)
	if err != nil {
		return nil, err
	}
	changeTopicTime(topics...)

	return topics, nil
}
