package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) GetTopicPurchase(ctx context.Context, uid, tid string) (*schema.TopicPurchase, error) {
	var topicPurchase schema.TopicPurchase
	if uid == "" || tid == "" {
		return nil, errors.New("uid and tid is empty")
	}

	find := bson.M{"uid": uid, "tid": tid}
	err := s.topicPurchaseColl.Find(ctx, find).One(&topicPurchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}

	return &topicPurchase, nil
}

func (s *MongoService) CreateTopicPurchase(ctx context.Context, uid, tid string) error {
	if uid == "" || tid == "" {
		return errors.New("uid and tid is empty")
	}

	tmNow := time.Now()

	topicPurchase := &schema.TopicPurchase{
		Uid:       uid,
		Tid:       tid,
		YesAmount: "0",
		NoAmount:  "0",
		CreatedAt: tmNow,
		UpdatedAt: tmNow,
	}
	_, err := s.topicPurchaseColl.InsertOne(ctx, topicPurchase)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrItemExist
		}
		return err
	}

	return nil
}

func (s *MongoService) QueryTopicPurchase(ctx context.Context, uid, tid string) ([]*schema.TopicPurchase, error) {
	var topicPurchases []*schema.TopicPurchase
	find := bson.M{}
	if uid != "" {
		find["uid"] = uid
	}
	if tid != "" {
		find["tid"] = tid
	}
	err := s.topicColl.Find(ctx, find).All(&topicPurchases)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return topicPurchases, nil
}

func (s *MongoService) UpsertTopicPurchase(ctx context.Context, topicPurchase *schema.TopicPurchase) error {
	topicPurchase.UpdatedAt = time.Now()
	filter := bson.M{"uid": topicPurchase.Uid, "tid": topicPurchase.Tid}
	_, err := s.topicPurchaseColl.Upsert(ctx, filter, topicPurchase)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoService) UpdateTopicPurchase(ctx context.Context, topicPurchase *schema.TopicPurchase) error {
	topicPurchase.UpdatedAt = time.Now()
	filter := bson.M{"uid": topicPurchase.Uid, "tid": topicPurchase.Tid}
	err := s.topicPurchaseColl.UpdateOne(ctx, filter, bson.M{"$set": topicPurchase})
	if err != nil {
		return err
	}

	return nil
}
