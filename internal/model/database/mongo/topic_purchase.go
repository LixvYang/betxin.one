package mongo

import (
	"context"
	"errors"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTopicPurchaseExist  error = errors.New("topic_purchase already exist")
	ErrNoSuchTopicPurchase error = errors.New("topic_purchase no exist")
)

func (s *MongoService) GetTopicPurchase(ctx context.Context, logger *zerolog.Logger, uid, tid string) (*schema.TopicPurchase, error) {
	var topicPurchase schema.TopicPurchase
	if uid == "" || tid == "" {
		return nil, errors.New("uid and tid is empty")
	}

	find := bson.M{"uid": uid, "tid": tid}

	err := s.topicColl.Find(ctx, find).One(&topicPurchase)
	if err == mongo.ErrNoDocuments {
		logger.Error().Err(err).Msg("mongo: not fount topic purchase")
		return nil, ErrNoSuchRefund
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error().Err(err).Msg("mongo: not fount topic purchase")
			return nil, ErrNoSuchRefund
		}
		logger.Error().Str("tid", tid).Err(err).Msg("mongo: get topic failed")
		return nil, err
	}

	return &topicPurchase, nil
}

func (s *MongoService) CreateTopicPurchase(ctx context.Context, logger *zerolog.Logger, topicPurchase *schema.TopicPurchase) error {
	_, err := s.topicColl.InsertOne(ctx, topicPurchase)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrTopicPurchaseExist
		}
		logger.Error().Err(err).Msg("mongo: create topic purchase failed")
		return err
	}

	return nil
}

func (s *MongoService) QueryTopicPurchase(ctx context.Context, logger *zerolog.Logger, uid, tid string) ([]*schema.TopicPurchase, error) {
	var topicPurchases []*schema.TopicPurchase
	find := bson.M{}
	if uid != "" {
		find["uid"] = uid
	}
	if tid != "" {
		find["tid"] = tid
	}
	err := s.topicColl.Find(ctx, find).All(&topicPurchases)
	if err == mongo.ErrNoDocuments {
		logger.Error().Err(err).Msg("mongo: not fount topic purchase")
		return nil, ErrNoSuchRefund
	}
	return topicPurchases, nil
}
