package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

func (s *MongoService) ListTopicByCid(c context.Context, logger *zerolog.Logger, cid int64, preId int64, pageSize int64) ([]*schema.Topic, error) {
	
	return nil, nil
}
