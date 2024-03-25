package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) ListCollects(ctx context.Context, uid string) ([]*schema.Collect, error) {
	var collects []*schema.Collect
	filter := bson.M{"uid": uid, "status": true}

	find := s.collectColl.Find(ctx, filter)

	err := find.All(&collects)
	if err != nil {
		return nil, err
	}
	return collects, nil
}

func (s *MongoService) GetCollectsByUid(ctx context.Context, uid string, limit, offset int64) ([]*schema.Collect, int64, error) {
	var collects []*schema.Collect
	filter := bson.M{"uid": uid, "status": true}

	find := s.collectColl.Find(ctx, filter)
	total, err := find.Count()
	if err != nil {
		return nil, 0, err
	}

	err = find.Limit(limit).Skip(offset).All(&collects)
	if err != nil {
		return nil, 0, err
	}

	return nil, total, nil
}

func (s *MongoService) UpsertCollect(ctx context.Context, uid, tid string, req *schema.Collect) error {
	_, err := s.collectColl.Upsert(ctx, bson.M{"uid": uid, "tid": tid}, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}

func (s *MongoService) GetCollectByUidTid(ctx context.Context, uid, tid string) (*schema.Collect, error) {
	var collect schema.Collect
	filter := bson.M{"uid": uid, "tid": tid}
	err := s.collectColl.Find(ctx, filter).One(&collect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return &collect, nil
}
