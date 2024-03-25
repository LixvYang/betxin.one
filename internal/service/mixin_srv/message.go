package mixin_srv

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	// 购买话题消息格式
	buyTopicMsgProto = `
🤖🤖🤖话题购买通知

ℹ️话题ID: %s
🗯️话题介绍: %s
🧩话题总数额: 
  - YES: %s
  - NO: %s

💸您的购买数额:
  - YES: %s
  - NO: %s

若您仍需要参与话题，请点击下方链接继续参与:
%s

⏲️话题结束时间(东八区): %s
`

	// 话题结束消息格式
	endTopicMsgProto = `
🤖🤖🤖话题结束通知

ℹ️话题ID: %s
🗯️话题介绍: %s
🧩话题总数额:
  - YES: %s
  - NO: %s
🟣话题结果: %s

💸您的购买数额:
  - YES: %s
  - NO: %s

🤑您将获得: %s
	
话题已结束，感谢您的参与！
详细信息: %s
`

	// 退款消息格式
	refundTopicMsgProto = `
🤖🤖🤖话题退款通知

ℹ️话题ID: %s
🗯️话题介绍: %s
✅退款行为: %s
💰退款数额: %s

💸您的购买数额仍有: 
  - YES: %s
  - NO: %s

若您仍需要参与话题，请点击下方链接继续参与:
%s

⏲️话题结束时间(东八区): %s
`

	// 话题已结束格式
	endedTopicMsgProto = `
😭😭😭话题已结束

ℹ️话题ID: %s
🗯️话题介绍: %s
🧩话题总数额: 
  - YES: %s
  - NO: %s

💸您的购买数额:
  - YES: %s
  - NO: %s

若您仍需要参与话题，请点击下方链接继续参与:
%s

⏲️话题结束时间(东八区): %s
`

	// TODO 距离话题截止退款时间 ...
	// TODO 距离话题结束时间...
)

func (m *MixinSrv) SendBuyTopicMessage(ctx context.Context, topic *schema.Topic, topicPurchase *schema.TopicPurchase) error {
	receiptId := topicPurchase.Uid
	topicLink := fmt.Sprintf("%s/topic/%s", m.appConf.WebUrl, topic.Tid)
	text := fmt.Sprintf(buyTopicMsgProto,
		topic.Tid,
		topic.Intro,
		topic.YesAmount,
		topic.NoAmount,
		topicPurchase.YesAmount,
		topicPurchase.NoAmount,
		topicLink,
		topic.EndTime.Format("2006-01-02 15:04:05"),
	)

	return m.SendEncryptedMessageWithRetry(ctx, receiptId, text)
}

// topicPurchase 需要传退款后的信息
func (m *MixinSrv) SendRefundTopicMessage(ctx context.Context, topic *schema.Topic, refund *schema.TopicRefundAction, topicPurchase *schema.TopicPurchase) error {
	receiptId := topicPurchase.Uid
	topicLink := fmt.Sprintf("%s/topic/%s", m.appConf.WebUrl, topic.Tid)
	var action string
	if refund.Action {
		action = "YES"
	} else {
		action = "NO"
	}
	text := fmt.Sprintf(refundTopicMsgProto,
		topic.Tid,
		topic.Intro,
		action,
		refund.Amount,
		topicPurchase.YesAmount,
		topicPurchase.NoAmount,
		topicLink,
		topic.EndTime.Format("2006-01-02 15:04:05"),
	)

	return m.SendEncryptedMessageWithRetry(ctx, receiptId, text)
}

// 批量发送话题结束消息
func (m *MixinSrv) SendStopTopicMessageMany(ctx context.Context, topic *schema.Topic, topicPurchaseRatioResps []*schema.TopicPurchaseRatio, result string) error {
	topicLink := fmt.Sprintf("%s/topic/%s", m.appConf.WebUrl, topic.Tid)
	var err error
	lo.ForEach(topicPurchaseRatioResps, func(item *schema.TopicPurchaseRatio, index int) {
		receiptId := item.Uid
		text := fmt.Sprintf(endTopicMsgProto,
			topic.Tid,
			topic.Intro,
			topic.YesAmount,
			topic.NoAmount,
			result,
			item.YesAmount,
			item.NoAmount,
			item.WinAmount,
			topicLink,
		)
		err = m.SendEncryptedMessageWithRetry(ctx, receiptId, text)
		if err != nil {
			return
		}
	})
	return nil
}

func (m *MixinSrv) SendStopTopicMessage(ctx context.Context, topic *schema.Topic, topicPurchaseRatioResp *schema.TopicPurchaseRatio, result string) error {
	receiptId := topicPurchaseRatioResp.Uid
	topicLink := fmt.Sprintf("%s/topic/%s", m.appConf.WebUrl, topic.Tid)
	text := fmt.Sprintf(endTopicMsgProto,
		topic.Tid,
		topic.Intro,
		topic.YesAmount,
		topic.NoAmount,
		result,
		topicPurchaseRatioResp.YesAmount,
		topicPurchaseRatioResp.NoAmount,
		topicPurchaseRatioResp.WinAmount,
		topicLink,
	)

	return m.SendEncryptedMessageWithRetry(ctx, receiptId, text)
}

const (
	defaultMaxMixinRetry = 3
)

func (m *MixinSrv) SendEncryptedMessageWithRetry(ctx context.Context, receiptId string, text string) error {
	var err error
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if err = m.SendEncryptedMessage(ctx, receiptId, text); err != nil {
			log.Error().Err(err).Msg("send message failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinSrv) SendEncryptedMessage(ctx context.Context, receiptId string, text string) error {
	sessions, err := m.Client.FetchSessions(ctx, []string{receiptId})
	if err != nil {
		log.Error().Err(err).Msg("fetch session failed")
		return err
	}

	_ = sessions

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(m.Client.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(text)),
	}

	if err := m.Client.EncryptMessageRequest(req, sessions); err != nil {
		log.Error().Err(err).Msg("EncryptMessageRequest failed")
		return err
	}

	_, err = m.Client.SendEncryptedMessages(ctx, []*mixin.MessageRequest{req})
	if err != nil {
		log.Error().Err(err).Msg("SendEncryptedMessages failed")
		return err
	}
	return nil
}
