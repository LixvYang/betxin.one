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
	ErrUserExist  error = errors.New("user already exist")
	ErrNoSuchUser error = errors.New("user no exist")
)

func (s *MongoService) CreateUser(ctx context.Context, log *zerolog.Logger, user *schema.User) (err error) {
	_, err = s.userColl.InsertOne(ctx, user)
	if err != nil {
		if isMongoDupeKeyError(err) {
			return ErrUserExist
		}
		log.Error().Msgf("failed to create user, error %v", err)
		return err
	}
	return nil
}

func (s *MongoService) GetUserByUid(ctx context.Context, log *zerolog.Logger, uid string) (*schema.User, error) {
	var user schema.User
	err := s.userColl.Find(ctx, bson.M{"uid": uid}).One(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchUser
		}
		log.Error().Msgf("failed to get user by uid, error %v", err)
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.Local()
	user.UpdatedAt = user.UpdatedAt.Local()
	return &user, nil
}

func (s *MongoService) UpdateUser(ctx context.Context, log *zerolog.Logger, uid string, user *schema.User) error {
	err := s.userColl.UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"$set": user})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoSuchUser
		}
		log.Error().Msgf("failed to update user, error %v", err)
		return err
	}
	return nil
}
