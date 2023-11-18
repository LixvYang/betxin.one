package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (th *TopicHandler) Delete(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	tid, err := th.checkDelete(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][checkDelete] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = th.storage.DeleteTopic(c, logger, tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][storage.DeleteTopic] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Int64("tid", tid).Msg("[Delete][storage.DeleteTopic]")
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}

func (th *TopicHandler) checkDelete(c *gin.Context) (int64, error) {
	tid, err := convert.StrToInt64(c.Param("tid"))
	if err != nil || tid == 0 {
		return 0, errors.New("[checkUpdateTopicInfoReq][StrToInt64] tid invalid")
	}
	return tid, nil
}
