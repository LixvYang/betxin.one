package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

func (s *MongoService) ListTopicByCid(c context.Context, logger *zerolog.Logger, cid int64, preId int64, pageSize int64) ([]*schema.Topic, error) {

	return nil, nil
}

func (s *MongoService) StopTopic(context.Context, *zerolog.Logger, string) error {
	return nil
}

func (s *MongoService) GetTopicsByCid(context.Context, *zerolog.Logger, int64) ([]*schema.Topic, error) {
	return nil, nil
}

func (s *MongoService) GetTopicByTid(context.Context, *zerolog.Logger, string) (*schema.Topic, error) {
	return nil, nil
}

func (s *MongoService) CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error {
	return nil
}

func (s *MongoService) DeleteTopic(context.Context, *zerolog.Logger, string) error {
	return nil
}

func (s *MongoService) UpdateTopic(context.Context, *zerolog.Logger, string, *schema.Topic) error {
	return nil
}

func (s *MongoService) GetTopicsByTids(ctx context.Context, logger *zerolog.Logger, tids []string) ([]*schema.Topic, error) {
	return nil, nil
}
