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
	ErrCategoryExist  error = errors.New("category already exist")
	ErrNoSuchCategory error = errors.New("category no exist")
)

func (s *MongoService) ListCategories(c context.Context, logger *zerolog.Logger) ([]*schema.Category, error) {
	var categorys []*schema.Category
	if err := s.categoryColl.Find(c, bson.M{}).All(&categorys); err != nil {
		logger.Error().Err(err).Msg("mongo: list categories failed")
		return nil, err
	}
	return nil, nil
}

func (s *MongoService) UpdateCategory(c context.Context, logger *zerolog.Logger, id int64, name string) error {
	err := s.categoryColl.UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCategory
		}
		logger.Error().Err(err).Msg("mongo: update category failed")
	}
	return nil
}

func (s *MongoService) CreateCategory(ctx context.Context, logger *zerolog.Logger, name string) error {
	_, err := s.categoryColl.InsertOne(ctx, bson.M{"name": name})
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrCategoryExist
		}
		logger.Error().Err(err).Msg("mongo: create category failed")
	}
	return nil
}
