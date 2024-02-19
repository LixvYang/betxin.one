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
	ErrBonuseExist  error = errors.New("bonuse already exist")
	ErrNoSuchBonuse error = errors.New("bonuse no exist")
)

func (s *MongoService) CreateBonuse(ctx context.Context, logger *zerolog.Logger, bonuse *schema.Bonuse) error {
	_, err := s.bonuseColl.InsertOne(ctx, bonuse)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrCategoryExist
		}
		logger.Error().Err(err).Msg("mongo: create category failed")
		return err
	}
	return nil
}

func (s *MongoService) GetBonuseByTraceId(ctx context.Context, logger *zerolog.Logger, traceId string) (*schema.Bonuse, error) {
	var bonuse schema.Bonuse
	err := s.bonuseColl.Find(ctx, bson.M{"trace_id": traceId}).One(&bonuse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchBonuse
		}
		logger.Error().Err(err).Msg("mongo: get category failed")
		return nil, err
	}
	return &bonuse, nil
}

func (s *MongoService) QueryBonuses(ctx context.Context, logger *zerolog.Logger, uid, tid string) ([]*schema.Bonuse, error) {
	find := bson.M{}
	if tid != "" {
		find["tid"] = tid
	}
	if uid != "" {
		find["uid"] = uid
	}

	bonuses := make([]*schema.Bonuse, 0)
	err := s.bonuseColl.Find(ctx, find).All(&bonuses)
	if err != nil {
		logger.Error().Err(err).Msg("mongo: query category failed")
		return nil, err
	}
	return bonuses, nil
}
