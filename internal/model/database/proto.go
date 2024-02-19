package database

import (
	"context"
	"time"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog"
)

type IUser interface {
	UpdateUser(ctx context.Context, log *zerolog.Logger, uid string, user *schema.User) error
	GetUserByUid(ctx context.Context, log *zerolog.Logger, uid string) (*schema.User, error)
	CreateUser(ctx context.Context, log *zerolog.Logger, user *schema.User) (err error)
}

type ITopic interface {
	StopTopic(context.Context, *zerolog.Logger, string) error
	GetTopicsByCid(context.Context, *zerolog.Logger, int64) ([]*schema.Topic, error)
	GetTopicByTid(context.Context, *zerolog.Logger, string) (*schema.Topic, error)
	CreateTopic(context.Context, *zerolog.Logger, *schema.Topic) error
	DeleteTopic(context.Context, *zerolog.Logger, string) error
	UpdateTopic(context.Context, *zerolog.Logger, string, *schema.Topic) error
	ListTopicByCid(c context.Context, logger *zerolog.Logger, cid int64, createdAt time.Time, pageSize int64) ([]*schema.Topic, int64, error)
	GetTopicsByTids(ctx context.Context, logger *zerolog.Logger, tids []string) ([]*schema.Topic, error)
}

type ICategory interface {
	CreateCategory(ctx context.Context, logger *zerolog.Logger, category *schema.Category) error
	GetCategoryById(ctx context.Context, logger *zerolog.Logger, id int64) (*schema.Category, error)
	ListCategories(c context.Context, logger *zerolog.Logger) ([]*schema.Category, error)
	UpdateCategory(ctx context.Context, logger *zerolog.Logger, id int64, name string) error
	DeleteCategory(ctx context.Context, logger *zerolog.Logger, id int64) error
}

type IBonuse interface {
	CreateBonuse(ctx context.Context, logger *zerolog.Logger, bonuse *schema.Bonuse) error
	GetBonuseByTraceId(ctx context.Context, logger *zerolog.Logger, traceId string) (*schema.Bonuse, error)
	QueryBonuses(ctx context.Context, logger *zerolog.Logger, uid, tid string, limit, offset int64) ([]*schema.Bonuse, int64, error)
}

type ICollect interface {
	CreateCollect(ctx context.Context, logger *zerolog.Logger, collect *schema.Collect) error
	ListCollects(ctx context.Context, logger *zerolog.Logger) ([]*schema.Collect, error)
	GetCollectByUid(ctx context.Context, logger *zerolog.Logger, uid string) ([]*schema.Collect, error)
	UpdateCollect(ctx context.Context, logger *zerolog.Logger, uid string, tid int64, status bool) error
}

type IRefund interface {
	CreateRefund(ctx context.Context, logger *zerolog.Logger, refund *schema.Refund) error
	GetRefundByTraceId(ctx context.Context, logger *zerolog.Logger, tracdId string) (*schema.Refund, error)
	ListRefundsWithQuery(ctx context.Context, logger *zerolog.Logger, limit, offset int64, tid, uid string, createdAt time.Time) ([]*schema.Refund, int64, error)
}

type ITopicPurchaseHistory interface {
	CreateTopicPurchaseHistory(ctx context.Context, logger *zerolog.Logger, purchaseHistory *schema.TopicPurchaseHistory) error
}

type ITopicPurchase interface {
	GetTopicPurchase(ctx context.Context, logger *zerolog.Logger, uid, tid string) (*schema.TopicPurchase, error)
	CreateTopicPurchase(ctx context.Context, logger *zerolog.Logger, topicPurchase *schema.TopicPurchase) error
	QueryTopicPurchase(ctx context.Context, logger *zerolog.Logger, uid, tid string) ([]*schema.TopicPurchase, error)
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

type Database interface {
	IUser
	ICategory
	ITopic
	ICollect
	IRefund
	ITopicPurchase
	ITopicPurchaseHistory
	IBonuse
}

func New(logger *zerolog.Logger, conf *config.AppConfig) Database {
	switch conf.Driver {
	case "mongo":
		return mongo.NewMongoService(logger, conf)
	default:
		logger.Panic().Msgf("driver: %s no impl", conf.Driver)
		return nil
	}
}
