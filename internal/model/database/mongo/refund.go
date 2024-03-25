package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateRefund(ctx context.Context, refund *schema.Refund) error {
	_, err := s.collectColl.InsertOne(ctx, refund)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetRefundByRequestId(ctx context.Context, requestId string) (*schema.Refund, error) {
	var refund schema.Refund
	err := s.refundColl.Find(ctx, bson.M{"request_id": requestId}).One(&refund)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return &refund, nil
}

func (s *MongoService) ListRefundsWithQuery(ctx context.Context, limit, offset int64, tid, uid string, createdAt time.Time) ([]*schema.Refund, int64, error) {
	var refunds []*schema.Refund
	var total int64
	var err error

	order := []string{}
	if !createdAt.IsZero() {
		order = append(order, "-created_at")
	}

	filter := bson.M{}

	if uid != "" {
		filter["uid"] = uid
	}
	if tid != "" {
		filter["tid"] = tid
	}

	total, err = s.refundColl.Find(ctx, filter).Count()
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, errors.New("refund no exist")
	}

	err = s.refundColl.Find(ctx, filter).Sort(order...).Skip(offset).Limit(limit).All(&refunds)
	if err != nil {
		return nil, 0, err
	}
	return refunds, total, nil
}
