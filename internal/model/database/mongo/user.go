package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateUser(ctx context.Context, user *schema.User) (err error) {
	_, err = s.userColl.InsertOne(ctx, user)
	if err != nil {
		// mongo.IsDuplicateKeyError(err)
		if mongo.IsDuplicateKeyError(err) {
			return ErrItemExist
		}
		return err
	}
	return nil
}

func (s *MongoService) GetUserByUid(ctx context.Context, uid string) (*schema.User, error) {
	var user schema.User
	err := s.userColl.Find(ctx, bson.M{"uid": uid}).One(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchItem
		}
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.Local()
	user.UpdatedAt = user.UpdatedAt.Local()
	return &user, nil
}

func (s *MongoService) UpdateUser(ctx context.Context, uid string, user *schema.User) error {
	err := s.userColl.UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"$set": user})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchItem
		}
		return err
	}
	return nil
}
