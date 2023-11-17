package topic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/internal/model/cache"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/query"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/pkg/snowflake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Test(t *testing.T) {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		3306,
		"betxin",
	)
	err = snowflake.Init("2022-11-11", 1)
	assert.Nil(t, err)
	db, err = gorm.Open(mysql.Open(dns))
	assert.Nil(t, err)
	db = db.Debug()

	query.SetDefault(db)

	cache := cache.New(&configs.RedisConfig{
		Host:         "127.0.0.1",
		Password:     "123456",
		Port:         6379,
		DB:           0,
		PoolSize:     128,
		MinIdleConns: 100,
	})

	topicModel := NewTopicModel(query.Q, cache)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	schemaTopic := new(schema.Topic)
	schemaTopic.Cid = 1
	schemaTopic.Content = "Test Content Demo"
	schemaTopic.EndTime = time.Now().Add(3 * time.Hour).UnixMilli()
	schemaTopic.RefundEndTime = time.Now().Add(1 * time.Hour).UnixMilli()
	schemaTopic.Intro = "Intro demo"
	schemaTopic.Title = "Title demo"

	err = topicModel.CreateTopic(ctx, &log.Logger, schemaTopic)
	assert.Nil(t, err)

	sqlTopics, err := topicModel.GetTopicsByCid(ctx, &log.Logger, 1)
	assert.Nil(t, err)

	log.Info().Any("sqlTopics", sqlTopics).Send()

	err = topicModel.CheckTopicExist(ctx, &log.Logger, sqlTopics[0].Tid)
	assert.Nil(t, err)

	err = topicModel.CheckTopicStop(ctx, &log.Logger, sqlTopics[0].Tid)
	assert.Nil(t, err)

	sqlTopics[0].Content = "Test"
	err = topicModel.UpdateTopicInfo(ctx, &log.Logger, sqlTopics[0])
	assert.Nil(t, err)

	schemaTopicc, err := topicModel.GetTopicByTid(ctx, &log.Logger, sqlTopics[0].Tid)
	fmt.Printf("%#v", *schemaTopicc)
	assert.Nil(t, err)
	assert.Equal(t, sqlTopics[0].Content, schemaTopicc.Content)

	err = topicModel.StopTopic(ctx, &log.Logger, sqlTopics[0].Tid)
	assert.Nil(t, err)

	// TODO
	// topicModel.UpdateTopicTotalPrice()

	err = topicModel.DeleteTopic(ctx, &log.Logger, sqlTopics[0].Tid)
	assert.Nil(t, err)
}

func TestFind(t *testing.T) {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		3306,
		"betxin",
	)
	err = snowflake.Init("2022-11-11", 1)
	assert.Nil(t, err)
	db, err = gorm.Open(mysql.Open(dns))
	assert.Nil(t, err)
	db = db.Debug()

	query.SetDefault(db)

	cache := cache.New(&configs.RedisConfig{
		Host:         "127.0.0.1",
		Password:     "123456",
		Port:         6379,
		DB:           0,
		PoolSize:     128,
		MinIdleConns: 100,
	})
	topicModel := NewTopicModel(query.Q, cache)

	schemaTopicc, err := topicModel.GetTopicByTid(context.Background(), &log.Logger, 447204630376484864)
	assert.Nil(t, err)
	fmt.Printf("%#v", *schemaTopicc)
}
