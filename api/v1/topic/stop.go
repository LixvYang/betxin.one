package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (th *TopicHandler) Stop(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	tid, err := th.checkTid(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Stop][checkDelete] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = th.storage.StopTopic(c, logger, tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Stop][storage.StopTopic]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Msgf("[Stop] tid: %s stop", c.Param("tid"))
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}
