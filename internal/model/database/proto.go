package database

import (
	"context"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type Database interface {
	Close() error
	IUser
	ITopic
}

type IUser interface {
	CheckUser(context.Context, *zerolog.Logger, string) error
	GetUserByUid(context.Context, *zerolog.Logger, string) (*schema.User, error)
	CreateUser(context.Context, *zerolog.Logger, *schema.User) error
	DeleteUser(context.Context, *zerolog.Logger, string) error
	UpdateUser(context.Context, *zerolog.Logger, *schema.User) error
}

type ITopic interface {
	StopTopic(context.Context, *zerolog.Logger, int64) error
	CheckTopicExist(context.Context, *zerolog.Logger, int64) error
	CheckTopicStop(context.Context, *zerolog.Logger, int64) error
	GetTopicsByCid(context.Context, *zerolog.Logger, int64) ([]*schema.Topic, error)
	GetTopicByTid(context.Context, *zerolog.Logger, int64) (*schema.Topic, error)
	CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	DeleteTopic(context.Context, *zerolog.Logger, int64) error
	UpdateTopicInfo(context.Context, *zerolog.Logger, *schema.Topic) error
	ListTopicByCid(context.Context, *zerolog.Logger, int64, int64, int64) ([]*schema.Topic, error)
	// TODO 字段
	// UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error
	// SearchTopic(context.Context, *zerolog.Logger, ...any) ([]*schema.Topic, int, error)
	// ListTopics(context.Context, *zerolog.Logger) ([]*schema.Topic, int, error)
}

func New(conf *configs.AppConfig) Database {
	return mysql.NewMySqlService(conf)
}
