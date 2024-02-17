package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

func (s *MongoService) CreateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64) error {
	return nil
}

func (s *MongoService) ListCollects(ctx context.Context, logger *zerolog.Logger) ([]*schema.Collect, error) {
	return nil, nil
}

func (s *MongoService) GetCollectByUid(ctx context.Context, logger *zerolog.Logger, uid string) ([]*schema.Collect, error) {
	return nil, nil
}

func (s *MongoService) UpdateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64, status bool) (*schema.Collect, error) {
	return nil, nil
}
