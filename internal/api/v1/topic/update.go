package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (th *TopicHandler) UpdateTopicInfo(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	newTopic, err := th.checkUpdateTopicInfoReq(c)
	if err != nil {
		logger.Error().Str("tid", c.Param("tid")).Err(err).Msg("[UpdateTopicInfo][checkUpdateTopicInfoReq] err")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	err = th.topicStore.UpdateTopicInfo(c, newTopic)
	if err != nil {
		logger.Error().Str("tid", c.Param("tid")).Err(err).Msg("[UpdateTopicInfo][storage.UpdateTopicInfo] err")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	getTopicResp, err := th.getTopicResp(c, newTopic.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][storage.GetTopicByTid]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Any("getTopicResp", getTopicResp).Msg("[UpdateTopicInfo][storage.UpdateTopicInfo]")

	handler.SendResponse(c, errmsg.SUCCSE, getTopicResp)
}

func (th *TopicHandler) checkUpdateTopicInfoReq(c *gin.Context) (*core.Topic, error) {
	tid, err := th.checkTid(c)
	if err != nil {
		return nil, err
	}
	var req CreateTopicReq

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	newTopic, err := th.checkCreateReq(&req)
	if err != nil {
		return nil, err
	}
	newTopic.Tid = tid

	return newTopic, nil
}
