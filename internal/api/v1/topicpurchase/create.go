package topicpurchase

import (
	"errors"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateTopicPurchaseActionReq struct {
	RequestId string `json:"request_id" binding:"required"`
	Action    bool   `json:"action" binding:"required"` // 0 NO 1 YES
	Tid       string `json:"tid" binding:"required"`
	Uid       string `json:"uid" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"` // CNB Amount
	Memo      string `json:"memo" binding:"required"`
}

// 话题购买 行为表
func (tph *TopicPurchaseHandler) CreatePurchaseAction(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	req, err := tph.checkCreateTopicPurchaseReq(c, &logger)
	if err != nil {
		logger.Error().Err(err).Any("req", req).Msg("invalid request")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	err = tph.storage.CreateTopicPurchaseHistory(c, &schema.TopicPurchaseHistory{
		RequestID: req.RequestId,
		Tid:       req.Tid,
		Uid:       req.Uid,
		Action:    req.Action,
		Amount:    fmt.Sprintf("%d", req.Amount),
		Memo:      req.Memo,
		CreatedAt: time.Now(),
	})

	if err != nil {
		logger.Error().Err(err).Msg("create topic purchase history failed")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	handler.SendResponse(c, errmsg.SUCCES, nil)
}

func (tph *TopicPurchaseHandler) checkCreateTopicPurchaseReq(c *gin.Context, logger *zerolog.Logger) (*CreateTopicPurchaseActionReq, error) {
	var req CreateTopicPurchaseActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid request body")
	}

	if req.Amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	if req.Tid == "" {
		return nil, errors.New("invalid tid")
	}

	// 查找话题 若话题已删除 或话题已截止购买 则报错
	topic, err := tph.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		if err == mongo.ErrNoSuchItem {
			return nil, errors.New("topic not found")
		}
		return nil, err
	}

	if topic.IsStop {
		return nil, errors.New("topic is stop")
	}

	if topic.EndTime.Before(time.Now()) {
		return nil, errors.New("topic is end")
	}

	uid, ok := utils.GetUserId(c)
	if !ok {
		return nil, errors.New("invalid uid")
	}

	if req.Uid != uid {
		return nil, errors.New("invalid uid")
	}

	// uid 若没有 则报错
	if _, err := tph.storage.GetUserByUid(c, req.Uid); err != nil {
		return nil, errors.New("user not found")
	}

	return &req, nil
}
