package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (th *TopicHandler) Delete(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	tid, err := th.checkTid(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][checkDelete] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = th.topicSrv.DeleteTopic(c, &logger, tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][storage.DeleteTopic] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Str("tid", tid).Msg("[Delete][storage.DeleteTopic]")
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}

func (th *TopicHandler) checkTid(c *gin.Context) (string, error) {
	tid := c.Param("tid")
	_, err := uuid.FromString(tid)
	if err != nil {
		return "", err
	}
	return tid, nil
}
