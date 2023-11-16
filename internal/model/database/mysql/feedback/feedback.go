package feedback

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type FeedbackModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewFeedbackModel(query *query.Query, cache *cache.Cache) FeedbackModel {
	return FeedbackModel{
		db:    query,
		cache: cache,
	}
}

func (fm *FeedbackModel) CreateFeedback(*schema.Feedback) error {
	return nil
}

func (fm *FeedbackModel) ListFeedback(uid string) ([]*schema.Feedback, error) {
	return nil, nil
}

func (fm *FeedbackModel) UpdateFeedback(*schema.Feedback) error {
	return nil
}
func (fm *FeedbackModel) DeleteFeedback(uid, fid string) error {
	return nil
}
func (fm *FeedbackModel) GetFeedback(uid, fid string) (*schema.Feedback, error) {
	return nil, nil
}
