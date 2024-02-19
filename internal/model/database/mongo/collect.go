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

func (s *MongoService) CreateCollect(ctx context.Context, logger *zerolog.Logger, collect *schema.Collect) error {
	_, err := s.collectColl.InsertOne(ctx, collect)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrCollectExist
		}
		logger.Error().Err(err).Msg("mongo: create collect failed")
		return err
	}
	return nil
}

func (s *MongoService) ListCollects(ctx context.Context, logger *zerolog.Logger) ([]*schema.Collect, error) {
	return nil, nil
}

func (s *MongoService) GetCollectByUid(ctx context.Context, logger *zerolog.Logger, uid string) ([]*schema.Collect, error) {
	var collects []*schema.Collect
	err := s.collectColl.Find(ctx, bson.M{"uid": uid}).All(&collects)
	if err != nil {
		logger.Error().Err(err).Msg("mongo: list collects failed")
		return nil, err
	}
	return nil, nil
}

func (s *MongoService) UpdateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64, status bool) error {
	err := s.collectColl.UpdateOne(ctx, bson.M{"uid": uid, "tid": tid}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCollect
		}
		logger.Error().Err(err).Msg("mongo: update collect failed")
	}
	return nil
}
