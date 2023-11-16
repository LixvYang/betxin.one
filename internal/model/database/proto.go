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
	ICategoty
	ICollect
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
	ListTopicByCid(ctx context.Context, logger *zerolog.Logger, cid int64, preId int64, pageSize int64) ([]*schema.Topic, error)
	// TODO 字段
	// UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error
	// SearchTopic(context.Context, *zerolog.Logger, ...any) ([]*schema.Topic, int, error)
	// ListTopics(context.Context, *zerolog.Logger) ([]*schema.Topic, int, error)
}

type ICategoty interface {
	CheckCategory(name string) error
	CreateCategory(name string) error
	GetCategoryById(id int64) (*schema.Category, error)
	ListCategories() ([]*schema.Category, error)
	UpdateCategory(id int64, name string) error
	DeleteCategory(id int64)
}

type IBonuse interface {
	CreateBonuse(*schema.Bonuse) error
	GetBonuseByTraceId(string) (*schema.Bonuse, error)
	ListBonuses() ([]*schema.Bonuse, error)
	UpdateBonuse(*schema.Bonuse) error
	DeleteBonuse(string) error
	GetBonusesByUid(string) (*schema.Bonuse, error)
}

type ICollect interface {
	CheckCollect(uid, tid string) (*schema.Collect, error)
	CreateCollect(uid, tid string) error
	ListCollects() ([]*schema.Collect, error)
	GetCollectByUid(uid string) ([]*schema.Collect, error)
	UpdateCollect(uid, tid string, status bool) (*schema.Collect, error)
}

type IFeedback interface {
	CreateFeedback(*schema.Feedback) error
	ListFeedback(uid string) ([]*schema.Feedback, error)
	UpdateFeedback(*schema.Feedback) error
	DeleteFeedback(uid, fid string) error
	GetFeedback(uid, fid string) (*schema.Feedback, error)
}

type IMessage interface {
	CreateMessage(*schema.Message) error
	ListMessage(uid string) ([]*schema.Message, error)
	DeleteMessage(conversation_id string) error
	GetMessage(conversation_id string) (*schema.Message, error)
}

type ISnapshot interface {
	CreateSnapshot(*schema.Snapshot) error
	ListMessage(uid string) ([]*schema.Snapshot, error)
	DeleteSnapshot(trace_id string) error
	GetSnapshot(trace_id string) (*schema.Snapshot, error)
}

type IRefund interface {
	CreateRefund(*schema.Refund) error
	GetRefundByTraceId(string) (*schema.Refund, error)
	ListRefunds() ([]*schema.Refund, error)
	UpdateRefund()
	DeleteRefund(trace_id string) error
	GetRefundsByUid(uid string) ([]*schema.Refund, error)
}
type ITopicPurchase interface {
	CheckTopicPurchase(uid, tid string) error
	GetTopicPurchase(uid, tid string) (*schema.TopicPurchase, error)
	CreateTopicPurchase(*schema.TopicPurchase) error
}

func New(conf *configs.AppConfig) Database {
	return mysql.NewMySqlService(conf)
}
