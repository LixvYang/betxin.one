package refund

import (
	"errors"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/service/mixin_srv"
	"github.com/lixvyang/betxin.one/internal/utils"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type CreateTopicRefundReq struct {
	Tid    string `json:"tid" binding:"required"`
	Amount int64  `json:"amount" binding:"required"`
	Action bool   `json:"action" binding:"required"`
}

func (rh *RefundHandler) CreateTopicRefund(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	req, err := rh.checkCreateTopicRefundReq(c)
	if err != nil {
		logger.Error().Err(err).Msg("checkCreateTopicRefundReq error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	uid, ok := utils.GetUserId(c)
	if !ok {
		logger.Error().Err(err).Msg("get user id error")
		handler.SendResponse(c, errmsg.ERROR_OAUTH, nil)
		return
	}

	topicRefundRequestId := convert.NewUUID()
	topicRefundMemo := fmt.Sprintf("话题退款: %s", req.Tid)
	topicRefundAction := &schema.TopicRefundAction{
		RequestID: topicRefundRequestId,
		Uid:       uid,
		Tid:       req.Tid,
		Amount:    decimal.NewFromInt(req.Amount),
		Action:    req.Action,
		Memo:      topicRefundMemo,
	}

	err = rh.storage.HandleMixinTopicRefundAction(c, topicRefundAction)
	if err != nil {
		logger.Error().Err(err).Msg("HandleMixinTopicRefundAction error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	transferReq := &mixin_srv.TransferOneRequest{
		RequestId: topicRefundAction.RequestID,
		Member:    topicRefundAction.Uid,
		Amount:    topicRefundAction.Amount,
		Memo:      topicRefundMemo,
	}

	// TODO: 向用户转账 并且发送消息
	_, err = rh.mixinSrv.TransferOne(c, transferReq)
	if err != nil {
		bb, _ := convert.Marshal(topicRefundAction)
		_ = rh.cache.ZAdd(c, consts.TopicRefundFailedRequest, float64(time.Now().Unix()), string(bb))
		logger.Error().Err(err).Msg("mixinSrv.TransferOne error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	topic, err := rh.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("GetTopicByTid error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	topicPurcase, err := rh.storage.GetTopicPurchase(c, topicRefundAction.Uid, topicRefundAction.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("GetTopicByTid error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	err = rh.mixinSrv.SendRefundTopicMessage(c, topic, topicRefundAction, topicPurcase)
	if err != nil {
		logger.Error().Err(err).Msg("mixinSrv.SendRefundTopicMessage error")
	}

	handler.SendResponse(c, errmsg.SUCCES, nil)
}

/*
1. 退款
2. 获取uid
3. 判断话题停止退款时间
4. 获取uid的购买记录 判断请求退款数是否对的上，判断话题时间
5. 向用户转账 并且发送消息
*/
func (rh *RefundHandler) checkCreateTopicRefundReq(c *gin.Context) (*CreateTopicRefundReq, error) {
	now := time.Now()
	var req CreateTopicRefundReq
	uid, ok := utils.GetUserId(c)
	if !ok {
		return nil, errors.New("invaild user")
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if req.Amount <= 0 {
		return nil, errors.New("invaild Amount")
	}

	_, err := convert.VaildUUID(req.Tid)
	if err != nil {
		return nil, err
	}

	topic, err := rh.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		return nil, err
	}

	if topic.RefundEndTime.Before(now) {
		return nil, errors.New("topic refund time is over")
	}

	topicPurchase, err := rh.storage.GetTopicPurchase(c, req.Tid, uid)
	if err != nil {
		return nil, err
	}

	if req.Action {
		if decimal.NewFromInt(req.Amount).GreaterThan(convert.NewDecimalFromString(topicPurchase.YesAmount)) {
			return nil, errors.New("invaild Amount")
		}
	} else {
		if decimal.NewFromInt(req.Amount).GreaterThan(convert.NewDecimalFromString(topicPurchase.NoAmount)) {
			return nil, errors.New("invaild Amount")
		}
	}

	return nil, err
}
