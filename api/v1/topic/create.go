package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type CreateTopicReq struct {
	Cid           int64  `json:"cid"`
	Title         string `json:"title"`
	Intro         string `json:"intro"`
	Content       string `json:"content"`
	ImgURL        string `json:"img_url"`
	RefundEndTime int64  `json:"refund_end_time"`
	EndTime       int64  `json:"end_time"`
}

type CreateTopicResp struct {
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

	if err = t.checkCreateReq(&req); err != nil {
		logger.Error().Err(err).Msg("[Create][checkCreateReq] error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	schemaTopic := new(schema.Topic)
	copier.Copy(schemaTopic, &req)

	err = t.storage.CreateTopic(c, logger, schemaTopic)
	if err != nil {
		logger.Error().Err(err).Msg("[Create][CreateTopic] error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}
	logger.Info().Any("topic", schemaTopic).Msg("[Create][CreateTopic] info")
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}

func (t *TopicHandler) checkCreateReq(req *CreateTopicReq) error {
	if req.Cid < 0 {
		return errors.New("cid error")
	}

	_, ok := t.categoryMap[req.Cid]
	if !ok {
		return errors.New("cid error")
	}

	if req.EndTime < 0 || req.RefundEndTime < 0 {
		return errors.New("time error")
	}

	if req.EndTime < req.RefundEndTime {
		return errors.New("time error")
	}

	if len(req.Title) > 40 {
		return errors.New("title length too long")
	}

	return nil
}
