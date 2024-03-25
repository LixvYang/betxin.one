package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (th *TopicHandler) UpdateTopicInfo(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	tid := c.Param("tid")
	if tid == "" {
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	newTopic, err := th.checkUpdateTopicInfoReq(c, &logger)
	if err != nil {
		logger.Error().Str("tid", tid).Err(err).Msg("[UpdateTopicInfo][checkUpdateTopicInfoReq] err")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	err = th.storage.UpdateTopic(c, tid, newTopic)
	if err != nil {
		logger.Error().Str("tid", tid).Err(err).Msg("[UpdateTopicInfo][storage.UpdateTopicInfo] err")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	getTopicResp, err := th.getTopicResp(c, &logger, newTopic.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][storage.GetTopicByTid]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Any("getTopicResp", getTopicResp).Msg("[UpdateTopicInfo][storage.UpdateTopicInfo]")

	handler.SendResponse(c, errmsg.SUCCES, getTopicResp)
}

func (th *TopicHandler) checkUpdateTopicInfoReq(c *gin.Context, logger *zerolog.Logger) (*schema.Topic, error) {
	tid, err := th.checkTid(c)
	if err != nil {
		return nil, err
	}

	newTopic, err := th.checkCreateReq(c, logger)
	if err != nil {
		return nil, err
	}
	newTopic.Tid = tid

	return newTopic, nil
}
