package mysql

import (
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/bonuse"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/category"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/collect"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/feedback"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/message"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/refund"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/snapshot"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/topic"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/topicpurchase"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/user"
	"github.com/lixvyang/betxin.one/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewMySqlService(conf *configs.AppConfig) *MySQLService {
	m := new(MySQLService)
	if err := m.initDB(conf.MySQLConfig); err != nil {
		logger.Lg.Error().Err(err).Msg("[NewMySqlService][m.Init()] err")
		panic(err)
	}
	cache := cache.New(conf.RedisConfig)
	m.UserModel = user.NewUserModel(query.Q, cache)
	m.TopicModel = topic.NewTopicModel(query.Q, cache)
	m.BonuseModel = bonuse.NewBonuseModel(query.Q, cache)
	m.CollectModel = collect.NewCollectModel(query.Q, cache)
	m.CategoryModel = category.NewCategoryModel(query.Q, cache)
	m.MessageModel = message.NewMessageModel(query.Q, cache)
	m.SnapshotModel = snapshot.NewMessageModel(query.Q, cache)
	m.RefundModel = refund.NewMessageModel(query.Q, cache)
	m.TopicPurchaseModel = topicpurchase.NewTopicPurchaseModel(query.Q, cache)
	m.FeedbackModel = feedback.NewFeedbackModel(query.Q, cache)

	return m
}

type MySQLService struct {
	db *gorm.DB
	user.UserModel
	topic.TopicModel
	category.CategoryModel
	collect.CollectModel
	bonuse.BonuseModel
	message.MessageModel
	snapshot.SnapshotModel
	topicpurchase.TopicPurchaseModel
	refund.RefundModel
	feedback.FeedbackModel
}

func (m *MySQLService) initDB(conf *configs.MySQLConfig) error {

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DB,
	)

	// info, err := query.Use(m.db).Topic
	var err error
	m.db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm日志模式：Warn
		Logger: gLogger.Default.LogMode(gLogger.Warn),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Lg.Panic().Msgf("连接数据库失败,请检查参数: %+v", err)
		return err
	}

	sqlDB, _ := m.db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(5000)
	// SetConnMaxLifetiment 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour / 2)
	query.SetDefault(m.db)
	return nil
}

func (m *MySQLService) Close() error {
	return nil
}
