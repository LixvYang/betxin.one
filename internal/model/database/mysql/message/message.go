package message

import (
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
)

type MessageModel struct {
	db    *query.Query
	cache *cache.Cache
}

func NewMessageModel(query *query.Query, cache *cache.Cache) MessageModel {
	return MessageModel{
		db:    query,
		cache: cache,
	}
}

func (mm *MessageModel) CreateMessage(*schema.Message) error {
	return nil
}

func (mm *MessageModel) ListMessage(uid string) ([]*schema.Message, error) {
	return nil, nil
}

func (mm *MessageModel) DeleteMessage(conversation_id string) error {
	return nil
}

func (mm *MessageModel) GetMessage(conversation_id string) (*schema.Message, error) {
	return nil, nil
}
