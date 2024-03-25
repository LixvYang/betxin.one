package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog/log"
)

const (
	cronStopTopicInterval = 5 * time.Minute
)

// 定时停止话题
func (c *CronService) CronStopTopic() {
	// 每 5 分钟遍历所有未停止的话题，检查是否有超过停止时间的，如果有，则停止话题
	ctx := context.Background()

	ticker := time.NewTicker(cronStopTopicInterval)
	for range ticker.C {
		now := time.Now()
		topics, count, err := c.storage.ListAllTopics(ctx)
		if err != nil {
			log.Error().Err(err).Msg("cron stop topic failed")
			continue
		}

		if count == 0 {
			continue
		}

		for _, topic := range topics {
			if topic.EndTime.After(now) {
				break
			}

			// 该停止的话题
			_, err := c.storage.HandleTopicStopAction(ctx, &schema.TopicStopAction{
				Tid: topic.Tid,
			})
			if err != nil {
				c.mixinSrv.SendEncryptedMessageWithRetry(ctx, c.mixinSrv.User.App.CreatorID, fmt.Sprintf("停止话题tid: %s 出错", topic.Tid))
				log.Error().Err(err).Msg("cron stop topic failed")
				continue
			}

			// 发送消息给参与者
			log.Info().Str("topic_id", topic.Tid).Msgf("topic stopped")
			c.mixinSrv.SendEncryptedMessageWithRetry(ctx, c.mixinSrv.User.App.CreatorID, fmt.Sprintf("话题tid: %s 已停止, 请及时完成话题结算", topic.Tid))
		}
	}
}

// // 处理话题停止
// func HandleTopicStopAction(ctx context.Context, stopTopicInfo *schema.StopTopicActionResp) {
// 	topicTotalAmount := convert.NewDecimalFromString(stopTopicInfo.Topic.YesAmount).Add(convert.NewDecimalFromString(stopTopicInfo.Topic.NoAmount))

// 	if topicTotalAmount.IsZero() {
// 		return
// 	}

// 	if topicTotalAmount.LessThan(decimal.NewFromInt(100)) {
// 		return
// 	}

// 	// 计算每个参与者所占话题的比例
// 	// 平台运营方收取的手续费为 CNB 的 2%
// 	ownerAmount := topicTotalAmount.Mul(decimal.NewFromFloat(0.02))

// 	othersTotalAmount := topicTotalAmount.Sub(ownerAmount)

// 	// 计算每个人所占话题的比例, 并进行转账, 给参与者发放奖励，最终结果记录下来
// 	for i := 0; i < len(stopTopicInfo.TopicPurchases); i++ {

// 	}

// }
