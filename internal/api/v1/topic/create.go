package topic

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/utils"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/internal/utils/timeof"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type CreateTopicReq struct {
	Cid           int64  `json:"cid"`
	Title         string `json:"title"`
	Intro         string `json:"intro"`
	Content       string `json:"content"`
	ImgURL        string `json:"img_url"`
	RefundEndTime string `json:"refund_end_time"`
	EndTime       string `json:"end_time"`
}

func (t *TopicHandler) Create(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	var req CreateTopicReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][ShouldBindJSON] error")
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	createTopicArgs, err := t.checkCreateReq(&req)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][checkCreateReq] error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	err = t.topicStore.CreateTopic(c, createTopicArgs)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][CreateTopic] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Any("topic", createTopicArgs).Msg("[Create][CreateTopic] info")

	resp, err := t.topicStore.GetTopicByTid(c, createTopicArgs.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][GetTopicByTid] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	createTopicResp := new(GetTopicResp)
	copier.Copy(&createTopicResp, &resp)

	handler.SendResponse(c, errmsg.SUCCSE, createTopicResp)
}

func (t *TopicHandler) checkCreateReq(req *CreateTopicReq) (*core.Topic, error) {
	if req.Cid < 0 {
		return nil, errors.New("cid error")
	}

	_, ok := t.categoryMap[req.Cid]
	if !ok {
		return nil, errors.New("cid error")
	}

	endTime, ok := timeof.TimeOf(req.EndTime)
	if !ok {
		return nil, fmt.Errorf("timeof endTime err: %s", req.EndTime)
	}
	refundEndTime, _ := timeof.TimeOf(req.RefundEndTime)

	if endTime.Before(refundEndTime) {
		return nil, fmt.Errorf("endTime: %s before refundEndTime: %s", req.EndTime, req.EndTime)
	}

	if len(req.Title) > 40 {
		return nil, errors.New("title length too long")
	}

	var argv = new(core.Topic)
	copier.Copy(&argv, req)
	argv.EndTime = endTime
	argv.RefundEndTime = refundEndTime
	argv.Tid = utils.NewUUID()

	return argv, nil
}
