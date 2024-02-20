package topic

import (
	"errors"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/internal/utils/timeof"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	createTopicArgs, err := t.checkCreateReq(c, &logger)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][checkCreateReq] error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	err = t.topicSrv.CreateTopic(c, &logger, createTopicArgs)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][CreateTopic] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Any("topic", createTopicArgs).Msg("[Create][CreateTopic] info")

	resp, err := t.topicSrv.GetTopicByTid(c, &logger, createTopicArgs.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][GetTopicByTid] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	createTopicResp := new(GetTopicResp)
	copier.Copy(&createTopicResp, &resp)
	createTopicResp.Category, _ = t.categorySrv.GetCategoryById(c, &logger, resp.Cid)
	handler.SendResponse(c, errmsg.SUCCSE, createTopicResp)
}

func (t *TopicHandler) checkCreateReq(c *gin.Context, logger *zerolog.Logger) (*schema.Topic, error) {
	var req CreateTopicReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][ShouldBindJSON] error")
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return nil, err
	}

	if req.Cid < 0 {
		return nil, errors.New("cid error")
	}

	_, err = t.categorySrv.GetCategoryById(c, logger, req.Cid)
	if err != nil {
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

	var argv = new(schema.Topic)
	copier.Copy(&argv, req)
	argv.EndTime = endTime
	argv.RefundEndTime = refundEndTime
	argv.Tid = utils.NewUUID()
	argv.CreatedAt = time.Now()
	argv.UpdatedAt = time.Now()
	argv.YesCount = "0"
	argv.NoCount = "0"
	argv.NoCount = "0"
	argv.YesRatio = "0"
	argv.NoRatio = "0"

	return argv, nil
}
