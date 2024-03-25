package topic

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type StopTopicReq struct {
	Tid string `json:"tid" binding:"required"`
	Win string `json:"win" binding:"required,oneof=YES NO"`
}

func (th *TopicHandler) StopTopic(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	var req StopTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Err(err).Any("req", req).Msg("[Stop][ShouldBindJSON]")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	if req.Win != "YES" && req.Win != "NO" {
		logger.Error().Any("req", req).Msg("[Stop][ShouldBindJSON]")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	stopAction, err := th.storage.HandleTopicStopAction(c, &schema.TopicStopAction{Tid: req.Tid})
	if err != nil {
		logger.Error().Err(err).Msg("[Stop][storage.StopTopic]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = th.handleTopicStopAction(c, logger, stopAction, req.Win)
	if err != nil {
		logger.Error().Err(err).Msg("[Stop][storage.StopTopic]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	handler.SendResponse(c, errmsg.SUCCES, nil)
}

// 处理话题停止
func (th *TopicHandler) handleTopicStopAction(ctx context.Context, logger zerolog.Logger, stopTopicInfo *schema.StopTopicActionResp, win string) error {
	if win != "YES" && win != "NO" {
		return fmt.Errorf("invalid win value: %s", win)
	}

	topicTotalAmount := convert.NewDecimalFromString(stopTopicInfo.Topic.YesAmount).Add(convert.NewDecimalFromString(stopTopicInfo.Topic.NoAmount))
	if topicTotalAmount.IsZero() {
		return fmt.Errorf("topic total amount is zero")
	}

	if topicTotalAmount.LessThan(decimal.NewFromInt(1)) {
		return fmt.Errorf("topic total amount is less than 1")
	}

	// 平台运营方收取的手续费为 CNB 的 2%
	ownerAmount := topicTotalAmount.Mul(decimal.NewFromFloat(0.02))
	_, _ = th.mixinSrv.TransferOne(ctx, &mixin_srv.TransferOneRequest{
		Amount:    ownerAmount,
		Memo:      fmt.Sprintf("话题 <%s> 结束手续费", stopTopicInfo.Topic.Title),
		Member:    th.mixinSrv.User.App.CreatorID,
		RequestId: convert.NewUUID(),
	})

	othersTotalAmount := topicTotalAmount.Sub(ownerAmount)

	var userTopicPurchaseRatio []*schema.TopicPurchaseRatio
	if win == "YES" {
		lo.ForEach(stopTopicInfo.TopicPurchases, func(item *schema.TopicPurchase, index int) {
			if item.YesAmount == "0" {
				return
			}
			ratio := convert.DecimalDiv(item.YesAmount, othersTotalAmount.String())
			userTopicPurchaseRatio = append(userTopicPurchaseRatio, &schema.TopicPurchaseRatio{
				Uid:       item.Uid,
				WinRatio:  ratio.String(),
				WinAmount: ratio.Mul(othersTotalAmount),
				YesAmount: convert.NewDecimalFromString(item.YesAmount),
				NoAmount:  convert.NewDecimalFromString(item.NoAmount),
			})
		})
	} else {
		lo.ForEach(stopTopicInfo.TopicPurchases, func(item *schema.TopicPurchase, index int) {
			if item.NoAmount == "0" {
				return
			}

			ratio := convert.DecimalDiv(item.YesAmount, othersTotalAmount.String())
			userTopicPurchaseRatio = append(userTopicPurchaseRatio, &schema.TopicPurchaseRatio{
				Uid:       item.Uid,
				WinRatio:  ratio.String(),
				WinAmount: ratio.Mul(othersTotalAmount),
				YesAmount: convert.NewDecimalFromString(item.YesAmount),
				NoAmount:  convert.NewDecimalFromString(item.NoAmount),
			})
		})
	}

	var memberAmounts []mixin_srv.MemberAmount // 每个切片最多255个用户
	for i := 0; i < len(userTopicPurchaseRatio); i++ {
		if len(memberAmounts) == consts.MAX_UTXO_NUM {
			err := th.mixinSrv.TransferManyWithRetry(ctx, &mixin_srv.TransferManyRequest{
				RequestId:    convert.NewUUID(),
				Memo:         fmt.Sprintf("参与话题 <%s> 奖励", stopTopicInfo.Topic.Title),
				MemberAmount: memberAmounts,
			})
			if err != nil {
				logger.Error().Err(err).Msg("[Stop][mixinSrv.TransferManyWithRetry]")
				continue
			}
			memberAmounts = make([]mixin_srv.MemberAmount, 0)
		}

		memberAmounts = append(memberAmounts, mixin_srv.MemberAmount{
			Member: []string{userTopicPurchaseRatio[i].Uid},
			Amount: userTopicPurchaseRatio[i].WinAmount,
		})
	}

	if len(memberAmounts) > 0 {
		err := th.mixinSrv.TransferManyWithRetry(ctx, &mixin_srv.TransferManyRequest{
			RequestId:    convert.NewUUID(),
			Memo:         fmt.Sprintf("参与话题 <%s> 奖励", stopTopicInfo.Topic.Title),
			MemberAmount: memberAmounts,
		})
		if err != nil {
			logger.Error().Err(err).Msg("[Stop][mixinSrv.TransferManyWithRetry]")
		}
	}

	// 发送消息给用户
	err := th.mixinSrv.SendStopTopicMessageMany(ctx, &stopTopicInfo.Topic, userTopicPurchaseRatio, win)
	if err != nil {
		logger.Error().Err(err).Msg("[Stop][mixinSrv.SendStopTopicMessageMany]")
	}
	return nil
}
