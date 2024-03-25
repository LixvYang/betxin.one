package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) ListCategories(c context.Context) ([]*schema.Category, error) {
	var categorys []*schema.Category
	if err := s.categoryColl.Find(c, bson.M{}).All(&categorys); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *MongoService) UpdateCategory(ctx context.Context, id int64, name string) error {
	err := s.categoryColl.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
	}
	return nil
}

func (s *MongoService) CreateCategory(ctx context.Context, category *schema.Category) error {
	_, err := s.categoryColl.InsertOne(ctx, category)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetCategoryById(ctx context.Context, id int64) (*schema.Category, error) {
	var category schema.Category
	err := s.categoryColl.Find(ctx, bson.M{"_id": id}).One(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}
	return &category, nil
}

func (s *MongoService) DeleteCategory(ctx context.Context, id int64) error {
	err := s.categoryColl.Remove(ctx, bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
	}
	return nil
}

func (s *MongoService) upsertCategory(ctx context.Context, id int64, name string) error {
	_, err := s.categoryColl.Upsert(ctx, bson.M{"_id": id}, bson.M{"_id": id, "name": name})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}
