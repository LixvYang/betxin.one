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
	UpdateUser(ctx context.Context, uid string, user *schema.User) error
	GetUserByUid(ctx context.Context, uid string) (*schema.User, error)
	CreateUser(ctx context.Context, user *schema.User) (err error)
}

type ITopic interface {
	ListTopics(ctx context.Context, cid int64, limit, offset int64) ([]*schema.Topic, int64, error)
	StopTopic(context.Context, string) error
	GetTopicsByCid(context.Context, int64) ([]*schema.Topic, error)
	GetTopicByTid(context.Context, string) (*schema.Topic, error)
	CreateTopic(context.Context, *schema.Topic) error
	DeleteTopic(context.Context, string) error
	UpdateTopic(context.Context, string, *schema.Topic) error
	ListTopicByCid(c context.Context, cid int64, createdAt time.Time, pageSize int64) ([]*schema.Topic, int64, error)
	GetTopicsByTids(ctx context.Context, tids []string) ([]*schema.Topic, error)
}

type ICategory interface {
	CreateCategory(ctx context.Context, category *schema.Category) error
	GetCategoryById(ctx context.Context, id int64) (*schema.Category, error)
	ListCategories(c context.Context) ([]*schema.Category, error)
	UpdateCategory(ctx context.Context, id int64, name string) error
	DeleteCategory(ctx context.Context, id int64) error
}

type IBonuse interface {
	CreateBonuse(ctx context.Context, bonuse *schema.Bonuse) error
	GetBonuseByTraceId(ctx context.Context, traceId string) (*schema.Bonuse, error)
	QueryBonuses(ctx context.Context, uid, tid string, limit, offset int64) ([]*schema.Bonuse, int64, error)
}

type ICollect interface {
	GetCollectByUidTid(ctx context.Context, uid, tid string) (*schema.Collect, error)
	ListCollects(ctx context.Context, uid string) ([]*schema.Collect, error)
	GetCollectsByUid(ctx context.Context, uid string, limit, offset int64) ([]*schema.Collect, int64, error)
	UpsertCollect(ctx context.Context, uid, tid string, req *schema.Collect) error
}

type IRefund interface {
	CreateRefund(ctx context.Context, refund *schema.Refund) error
	GetRefundByTraceId(ctx context.Context, tracdId string) (*schema.Refund, error)
	ListRefundsWithQuery(ctx context.Context, limit, offset int64, tid, uid string, createdAt time.Time) ([]*schema.Refund, int64, error)
}

type ITopicPurchaseHistory interface {
	CreateTopicPurchaseHistory(ctx context.Context, purchaseHistory *schema.TopicPurchaseHistory) error
}

type ITopicPurchase interface {
	GetTopicPurchase(ctx context.Context, uid, tid string) (*schema.TopicPurchase, error)
	CreateTopicPurchase(ctx context.Context, topicPurchase *schema.TopicPurchase) error
	QueryTopicPurchase(ctx context.Context, uid, tid string) ([]*schema.TopicPurchase, error)
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

type ISnapshot interface {
	InsertSnapshot(ctx context.Context, snapshot *schema.Snapshot) error
	GetSnapshotByRequestId(ctx context.Context, requestId string) (*schema.Snapshot, error)
	GetLastestSnapshot(ctx context.Context) (*schema.Snapshot, error)
	GetSnapshotCount(ctx context.Context) (int64, error)
}

type Database interface {
	IUser
	ICategory
	ITopic
	ICollect
	IRefund
	ITopicPurchase
	ITopicPurchaseHistory
	IBonuse
	ISnapshot
}

func New(logger *zerolog.Logger, conf *config.AppConfig) *mongo.MongoService {
	switch conf.Driver {
	case "mongo":
		return mongo.NewMongoService(conf.MongoConfig)
	default:
		logger.Panic().Msgf("driver: %s no impl", conf.Driver)
		return nil
	}
}
