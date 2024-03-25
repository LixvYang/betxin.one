package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/stretchr/testify/assert"
)

func Test_HandleMixinTopicDepositAction(t *testing.T) {
	assert := assert.New(t)
	mongoService := NewMongoService(
		&config.MongoConfig{
			Host: "127.0.0.1",
			Port: 30001,
			DB:   "betxinonetest",
			// Username: "mongoadmin",
			// Password: "mongoadmin",
		},
	)
	assert.Nil(mongoService.client.Ping(10))

	ctx := context.Background()
	tid := convert.NewUUID()
	now := time.Now()
	testTopic := schema.Topic{
		// Tid:       "123",
		Tid:       tid,
		Title:     "test-topic",
		Intro:     "test-intro",
		Content:   "test-content",
		YesAmount: "15.123",
		NoAmount:  "222.133",
		CreatedAt: now,
	}

	err := mongoService.CreateTopic(ctx, &testTopic)
	assert.Nil(err)

	topicS, err := mongoService.GetTopicByTid(ctx, tid)
	assert.Nil(err)

	fmt.Printf("%+v", topicS)
	uid := convert.NewUUID()
	// 用户购买记录
	err = mongoService.HandleMixinTopicDepositAction(ctx, &schema.TopicBuyAction{
		RequestID: convert.NewUUID(),
		Tid:       tid,
		Uid:       uid,
		Amount:    convert.NewDecimalFromString("10.001"),
		Action:    false,
	})
	assert.Nil(err)

	// 验证用户购买记录
	userBuyRecord, err := mongoService.GetTopicPurchase(ctx, uid, tid)
	assert.Nil(err)
	assert.Equal(tid, userBuyRecord.Tid)
	assert.Equal(uid, userBuyRecord.Uid)
	assert.Equal(convert.NewDecimalFromString("10.001").String(), userBuyRecord.NoAmount)
	assert.Equal(convert.NewDecimalFromString("0").String(), userBuyRecord.YesAmount)

	reqId := convert.NewUUID()
	err = mongoService.HandleMixinTopicDepositAction(ctx, &schema.TopicBuyAction{
		RequestID: reqId,
		Tid:       tid,
		Uid:       uid,
		Amount:    convert.NewDecimalFromString("5.00000001"),
		Action:    true,
	})
	assert.Nil(err)

	userBuyRecord, err = mongoService.GetTopicPurchase(ctx, uid, tid)
	assert.Nil(err)
	assert.Equal(convert.NewDecimalFromString("10.001").String(), userBuyRecord.NoAmount)
	assert.Equal(convert.NewDecimalFromString("5.00000001").String(), userBuyRecord.YesAmount)

	topic, err := mongoService.GetTopicByTid(ctx, tid)
	assert.Nil(err)
	assert.Equal(convert.NewDecimalFromString("20.12300001").String(), topic.YesAmount)
	assert.Equal(convert.NewDecimalFromString("232.134").String(), topic.NoAmount)

	// 用户退款记录
	err = mongoService.HandleMixinTopicRefundAction(ctx, &schema.TopicRefundAction{
		RequestID: convert.NewUUID(),
		Tid:       tid,
		Uid:       uid,
		Amount:    convert.NewDecimalFromString("0.001"),
		Action:    true,
	})
	assert.Nil(err)
	topic, err = mongoService.GetTopicByTid(ctx, tid)
	assert.Nil(err)
	assert.Equal(convert.NewDecimalFromString("20.12200001").String(), topic.YesAmount)

	err = mongoService.HandleMixinTopicRefundAction(ctx, &schema.TopicRefundAction{
		RequestID: convert.NewUUID(),
		Tid:       tid,
		Uid:       uid,
		Amount:    convert.NewDecimalFromString("0.123"),
		Action:    false,
	})
	assert.Nil(err)
	topic, err = mongoService.GetTopicByTid(ctx, tid)
	assert.Nil(err)
	assert.Equal(convert.NewDecimalFromString("232.011").String(), topic.NoAmount)
}
