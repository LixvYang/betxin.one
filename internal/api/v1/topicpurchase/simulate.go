package topicpurchase

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type PurchaseSimulateReq struct {
	Tid       string `json:"tid" binding:"required"`
	WinAction bool   `json:"win_action" binding:"required"`
	Amount    string `json:"amount" binding:"required"`

	uid    string `json:"-"`
	amount decimal.Decimal
}

type PurchaseSimulateResp struct {
	YesAmount string `json:"yes_amount"`
	NoAmount  string `json:"no_amount"`
}

// 处理模拟购买接口
func (tph *TopicPurchaseHandler) Simulate(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	var resp PurchaseSimulateResp
	req, err := tph.checkPurchaseSimulateReq(c)
	if err != nil {
		logger.Error().Any("req", req).Err(err).Msg("check purchase simulate request error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	userTopicPurchase, err := tph.storage.GetTopicPurchase(c, req.uid, req.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("get user topic purchase error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	topic, err := tph.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR_TOPICS_NOT_FOUND, nil)
		return
	}

	topicTotalAmount := convert.NewDecimalFromString(topic.YesAmount).Add(convert.NewDecimalFromString(topic.NoAmount))
	if topicTotalAmount.IsZero() {
		if req.WinAction {
			resp.YesAmount = req.amount.Mul(decimal.NewFromFloat(0.98)).String()
		} else {
			resp.NoAmount = req.amount.Mul(decimal.NewFromFloat(0.98)).String()
		}
		handler.SendResponse(c, errmsg.SUCCES, resp)
		return
	}

	// othersTotalAmount := topicTotalAmount.Mul(decimal.NewFromFloat(0.98))
	if req.WinAction {
		userTopicPurchase.YesAmount = convert.NewDecimalFromString(userTopicPurchase.YesAmount).Add(req.amount).String()
	} else {
		userTopicPurchase.NoAmount = convert.NewDecimalFromString(userTopicPurchase.NoAmount).Add(req.amount).String()
	}
	resp.YesAmount = convert.NewDecimalFromString(userTopicPurchase.YesAmount).Div(topicTotalAmount).Mul(decimal.NewFromFloat(0.98)).String()
	resp.NoAmount = convert.NewDecimalFromString(userTopicPurchase.NoAmount).Div(topicTotalAmount).Mul(decimal.NewFromFloat(0.98)).String()
	handler.SendResponse(c, errmsg.SUCCES, resp)
}

func (tph *TopicPurchaseHandler) checkPurchaseSimulateReq(c *gin.Context) (*PurchaseSimulateReq, error) {
	var req PurchaseSimulateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	_, err := tph.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		return nil, err
	}

	uid, ok := utils.GetUserId(c)
	if !ok {
		return nil, errors.New("user id not found in request context")
	}
	req.uid = uid

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, err
	}

	req.amount = amount
	return nil, nil
}
