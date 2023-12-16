package mysql

import (
	"fmt"
	"time"

	configs "github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/category"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/topic"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/store/user"
	"github.com/rs/zerolog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewMySqlService(logger *zerolog.Logger, conf *configs.AppConfig) *MySQLService {
	m := new(MySQLService)
	if err := m.initDB(logger, conf.MySQLConfig); err != nil {
		logger.Error().Err(err).Msg("[NewMySqlService][m.Init()] err")
		panic(err)
	}

	h := store.MustInit(conf)
	m.UserStore = user.New(h)
	m.TopicStore = topic.New(h)
	m.CategoryStore = category.New(h)

	return m
}

type MySQLService struct {
	db *gorm.DB
	core.UserStore
	core.TopicStore
	core.CategoryStore
}

func (m *MySQLService) initDB(logger *zerolog.Logger, conf *configs.MySQLConfig) error {
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
		logger.Panic().Msgf("连接数据库失败,请检查参数: %+v", err)
		return err
	}

	sqlDB, _ := m.db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(1000)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(5000)
	// SetConnMaxLifetiment 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour / 2)
	return nil
}

func (m *MySQLService) Close() error {
	return nil
}
