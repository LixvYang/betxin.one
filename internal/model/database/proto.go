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

type ITopic interface{}
type IITopic interface {
	StopTopic(context.Context, *zerolog.Logger, string) error
	CheckTopicExist(context.Context, *zerolog.Logger, string) error
	CheckTopicStop(context.Context, *zerolog.Logger, string) error
	GetTopicTotalPrice(context.Context, *zerolog.Logger, string) (string, error)
	GetTopicsByCid(context.Context, *zerolog.Logger, string) ([]*schema.Topic, int, error)
	GetTopicByTid(context.Context, *zerolog.Logger, string) (*schema.Topic, error)
	CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	DeleteTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	UpdateTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	// TODO 字段
	UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error
	SearchTopic(context.Context, *zerolog.Logger, ...any)
	ListTopics(context.Context, *zerolog.Logger)
}

func New(conf *configs.AppConfig) Database {
	return mysql.NewMySqlService(conf)
}
