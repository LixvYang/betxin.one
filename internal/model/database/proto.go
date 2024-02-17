package database

import (
	"context"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type Database interface {
	IUser
	ICategory
	ITopic
	ICollect
}

type IUser interface {
	// CheckUser(context.Context, *zerolog.Logger, string) error
	// GetUserByUid(context.Context, *zerolog.Logger, string) (*schema.User, error)
	// CreateUser(context.Context, *zerolog.Logger, *schema.User) error
	// DeleteUser(context.Context, *zerolog.Logger, string) error
	UpdateUser(ctx context.Context, log *zerolog.Logger, uid string, user *schema.User) error
	GetUserByUid(ctx context.Context, log *zerolog.Logger, uid string) (*schema.User, error)
	CreateUser(ctx context.Context, log *zerolog.Logger, user *schema.User) (err error)
}

type ITopic interface {
	StopTopic(context.Context, *zerolog.Logger, string) error
	// CheckTopicExist(context.Context, *zerolog.Logger, int64) error
	// CheckTopicStop(context.Context, *zerolog.Logger, int64) error
	GetTopicsByCid(context.Context, *zerolog.Logger, int64) ([]*schema.Topic, error)
	GetTopicByTid(context.Context, *zerolog.Logger, string) (*schema.Topic, error)
	CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	DeleteTopic(context.Context, *zerolog.Logger, string) error
	UpdateTopic(context.Context, *zerolog.Logger, string, *schema.Topic) error
	ListTopicByCid(c context.Context, logger *zerolog.Logger, cid int64, preId int64, pageSize int64) ([]*schema.Topic, error)
	GetTopicsByTids(ctx context.Context, logger *zerolog.Logger, tids []string) ([]*schema.Topic, error)
	// TODO 字段
	// UpdateTopicTotalPrice(context.Context, *zerolog.Logger, *schema.Topic) error
	// SearchTopic(context.Context, *zerolog.Logger, ...any) ([]*schema.Topic, int, error)
	// ListTopics(context.Context, *zerolog.Logger) ([]*schema.Topic, int, error)
}

type ICategory interface {
	CreateCategory(ctx context.Context, logger *zerolog.Logger, category *schema.Category) error
	GetCategoryById(ctx context.Context, logger *zerolog.Logger, id int64) (*schema.Category, error)
	ListCategories(c context.Context, logger *zerolog.Logger) ([]*schema.Category, error)
	UpdateCategory(ctx context.Context, logger *zerolog.Logger, id int64, name string) error
	DeleteCategory(ctx context.Context, logger *zerolog.Logger, id int64) error
}

// type IBonuse interface {
// 	CreateBonuse(*schema.Bonuse) error
// 	GetBonuseByTraceId(string) (*schema.Bonuse, error)
// 	ListBonuses() ([]*schema.Bonuse, error)
// 	UpdateBonuse(*schema.Bonuse) error
// 	DeleteBonuse(string) error
// 	GetBonusesByUid(string) (*schema.Bonuse, error)
// }

type ICollect interface {
	// CheckCollect(uid string, tid int64) (*schema.Collect, error)
	CreateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64) error
	ListCollects(ctx context.Context, logger *zerolog.Logger) ([]*schema.Collect, error)
	GetCollectByUid(ctx context.Context, logger *zerolog.Logger, uid string) ([]*schema.Collect, error)
	UpdateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64, status bool) (*schema.Collect, error)
}

// type IFeedback interface {
// 	CreateFeedback(*schema.Feedback) error
// 	ListFeedback(uid string) ([]*schema.Feedback, error)
// 	UpdateFeedback(*schema.Feedback) error
// 	DeleteFeedback(uid, fid string) error
// 	GetFeedback(uid, fid string) (*schema.Feedback, error)
// }

// type IMessage interface {
// 	CreateMessage(*schema.Message) error
// 	ListMessage(uid string) ([]*schema.Message, error)
// 	DeleteMessage(conversation_id string) error
// 	GetMessage(conversation_id string) (*schema.Message, error)
// }

// type ISnapshot interface {
// 	CreateSnapshot(*schema.Snapshot) error
// 	ListMessage(uid string) ([]*schema.Snapshot, error)
// 	DeleteSnapshot(trace_id string) error
// 	GetSnapshot(trace_id string) (*schema.Snapshot, error)
// }

// type IRefund interface {
// 	CreateRefund(*schema.Refund) error
// 	GetRefundByTraceId(string) (*schema.Refund, error)
// 	ListRefunds() ([]*schema.Refund, error)
// 	UpdateRefund()
// 	DeleteRefund(trace_id string) error
// 	GetRefundsByUid(uid string) ([]*schema.Refund, error)
// }
// type ITopicPurchase interface {
// 	CheckTopicPurchase(uid, tid string) error
// 	GetTopicPurchase(uid, tid string) (*schema.TopicPurchase, error)
// 	CreateTopicPurchase(*schema.TopicPurchase) error
// }

func New(logger *zerolog.Logger, conf *config.AppConfig) Database {
	switch conf.Driver {
	case "mongo":
		return mongo.NewMongoService(logger, conf)
	default:
		logger.Panic().Msgf("driver: %s no impl", conf.Driver)
		return nil
	}
}
