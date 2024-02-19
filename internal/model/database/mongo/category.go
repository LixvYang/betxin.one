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

func (s *MongoService) UpdateCategory(ctx context.Context, logger *zerolog.Logger, id int64, name string) error {
	err := s.categoryColl.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCategory
		}
		logger.Error().Err(err).Msg("mongo: update category failed")
	}
	return nil
}

func (s *MongoService) CreateCategory(ctx context.Context, logger *zerolog.Logger, category *schema.Category) error {
	_, err := s.categoryColl.InsertOne(ctx, category)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrCategoryExist
		}
		logger.Error().Err(err).Msg("mongo: create category failed")
		return err
	}
	return nil
}

func (s *MongoService) GetCategoryById(ctx context.Context, logger *zerolog.Logger, id int64) (*schema.Category, error) {
	var category schema.Category
	err := s.categoryColl.Find(ctx, bson.M{"_id": id}).One(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchCategory
		}
		logger.Error().Err(err).Msg("mongo: get category failed")
		return nil, err
	}
	return &category, nil
}

func (s *MongoService) DeleteCategory(ctx context.Context, logger *zerolog.Logger, id int64) error {
	err := s.categoryColl.Remove(ctx, bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCategory
		}
		logger.Error().Err(err).Msg("mongo: delete category failed")
	}
	return nil
}

func (s *MongoService) upsertCategory(ctx context.Context, logger *zerolog.Logger, id int64, name string) error {
	_, err := s.categoryColl.Upsert(ctx, bson.M{"_id": id}, bson.M{"_id": id, "name": name})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchCategory
		}
		logger.Error().Err(err).Msg("mongo: update category failed")
	}
	return nil
}
