package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTopicExist  error = errors.New("topic already exist")
	ErrNoSuchTopic error = errors.New("topic no exist")
)

func (s *MongoService) ListTopicByCid(ctx context.Context, logger *zerolog.Logger, cid int64, createdAt time.Time, pageSize int64) ([]*schema.Topic, int64, error) {
	var topics []*schema.Topic
	var total int64
	var err error

	filter := bson.M{"created_at": bson.M{"$lte": createdAt}, "cid": cid}
	find := s.topicColl.Find(ctx, filter)
	total, err = find.Count()
	if err != nil {
		return nil, total, err
	}

	if total == 0 {
		return nil, total, ErrNoSuchTopic
	}
	err = find.Limit(pageSize).Sort("-created_at").All(&topics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, total, ErrNoSuchTopic
		}
		return nil, total, err
	}
	return topics, total, nil
}

func (s *MongoService) StopTopic(ctx context.Context, logger *zerolog.Logger, tid string) error {
	err := s.topicColl.UpdateOne(ctx, bson.M{"tid": tid}, bson.M{"$set": bson.M{"is_stop": true}})
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrTopicExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetTopicsByCid(ctx context.Context, logger *zerolog.Logger, cid int64) ([]*schema.Topic, error) {
	var topics []*schema.Topic
	err := s.topicColl.Find(ctx, bson.M{"cid": cid}).All(&topics)
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func (s *MongoService) GetTopicByTid(ctx context.Context, logger *zerolog.Logger, tid string) (*schema.Topic, error) {
	var topic *schema.Topic
	err := s.topicColl.Find(ctx, bson.M{"tid": tid}).One(&topic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchTopic
		}
		return nil, err
	}
	return topic, nil
}

func (s *MongoService) CreateTopic(ctx context.Context, logger *zerolog.Logger, topic *schema.Topic) error {
	_, err := s.topicColl.InsertOne(ctx, topic)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrTopicExist
		}
		return err
	}
	return nil
}

func (s *MongoService) DeleteTopic(ctx context.Context, logger *zerolog.Logger, tid string) error {
	err := s.topicColl.Remove(ctx, bson.M{"tid": tid})
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrTopicExist
		}
		return err
	}
	return nil
}

func (s *MongoService) UpdateTopic(ctx context.Context, logger *zerolog.Logger, tid string, topic *schema.Topic) error {
	err := s.topicColl.UpdateOne(ctx, bson.M{"tid": tid}, bson.M{"$set": topic})
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrTopicExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetTopicsByTids(ctx context.Context, logger *zerolog.Logger, tids []string) ([]*schema.Topic, error) {
	var topics []*schema.Topic
	err := s.topicColl.Find(ctx, bson.M{"tid": bson.M{"$in": tids}}).All(&topics)
	if err != nil {
		return nil, err
	}
	return topics, nil
}
