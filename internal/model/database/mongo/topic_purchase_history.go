package mongo

import (
	"context"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

func (s *MongoService) CreateTopicPurchaseHistory(ctx context.Context, logger *zerolog.Logger, purchaseHistory *schema.TopicPurchaseHistory) error {
	_, err := s.topicPurchaseHistoryColl.InsertOne(ctx, purchaseHistory)
	if err != nil {
		logger.Error().Err(err).Msg("create topic purchase history failed")
		return err
	}
	return nil
}
