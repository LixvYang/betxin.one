package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) QueryTopicPurchaseHistory(ctx context.Context, requestId string) (*schema.TopicPurchaseHistory, error) {
	history := &schema.TopicPurchaseHistory{}
	err := s.topicPurchaseHistoryColl.Find(ctx, bson.M{"request_id": requestId}).One(&history)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return history, nil
}

func (s *MongoService) CreateTopicPurchaseHistory(ctx context.Context, purchaseHistory *schema.TopicPurchaseHistory) error {
	_, err := s.topicPurchaseHistoryColl.InsertOne(ctx, purchaseHistory)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) FinishedTopicPurchaseHistory(ctx context.Context, requestId string) error {
	filter := bson.M{"request_id": requestId}
	update := bson.M{"$set": bson.M{"finished": true}}
	err := s.topicPurchaseHistoryColl.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}

func (s *MongoService) UpsertTopicPurchaseHistory(ctx context.Context, purchaseHistory *schema.TopicPurchaseHistory) error {
	filter := bson.M{"request_id": purchaseHistory.RequestID}
	_, err := s.topicPurchaseHistoryColl.Upsert(ctx, filter, purchaseHistory)
	if err != nil {
		return err
	}
	return nil
}
