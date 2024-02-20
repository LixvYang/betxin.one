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
	ErrCollectExist  error = errors.New("collect already exist")
	ErrNoSuchCollect error = errors.New("collect no exist")
)

func (s *MongoService) ListCollects(ctx context.Context, logger *zerolog.Logger, uid string) ([]*schema.Collect, error) {
	var collects []*schema.Collect
	filter := bson.M{"uid": uid, "status": true}

	find := s.collectColl.Find(ctx, filter)

	err := find.All(&collects)
	if err != nil {
		return nil, err
	}
	return collects, nil
}

func (s *MongoService) GetCollectsByUid(ctx context.Context, logger *zerolog.Logger, uid string, limit, offset int64) ([]*schema.Collect, int64, error) {
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

func (s *MongoService) UpsertCollect(ctx context.Context, logger *zerolog.Logger, uid, tid string, req *schema.Collect) error {
	_, err := s.collectColl.Upsert(ctx, bson.M{"uid": uid, "tid": tid}, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCollect
		}
		return err
	}
	return nil
}

func (s *MongoService) GetCollectByUidTid(ctx context.Context, logger *zerolog.Logger, uid, tid string) (*schema.Collect, error) {
	var collect schema.Collect
	filter := bson.M{"uid": uid, "tid": tid}
	err := s.collectColl.Find(ctx, filter).One(&collect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchCollect
		}
		return nil, err
	}
	return &collect, nil
}
