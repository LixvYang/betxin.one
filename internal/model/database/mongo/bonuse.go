package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateBonuse(ctx context.Context, bonuse *schema.Bonuse) error {
	_, err := s.bonuseColl.InsertOne(ctx, bonuse)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetBonuseByTraceId(ctx context.Context, traceId string) (*schema.Bonuse, error) {
	var bonuse schema.Bonuse
	err := s.bonuseColl.Find(ctx, bson.M{"trace_id": traceId}).One(&bonuse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return &bonuse, nil
}

func (s *MongoService) QueryBonuses(ctx context.Context, uid, tid string, limit, offset int64) ([]*schema.Bonuse, int64, error) {
	var err error
	var bonuses []*schema.Bonuse
	var total int64
	find := bson.M{}
	if tid != "" {
		find["tid"] = tid
	}
	if uid != "" {
		find["uid"] = uid
	}

	total, err = s.bonuseColl.Find(ctx, find).Count()
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, ErrNoSuchItem
	}

	err = s.bonuseColl.Find(ctx, find).Sort("-created_at").Skip(offset).Limit(limit).All(&bonuses)
	if err != nil {
		return nil, 0, err
	}

	return bonuses, total, nil
}
