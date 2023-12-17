package collect

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (ch *CollectHandler) Delete(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	uid, ok := c.MustGet("uid").(string)
	if !ok {
		logger.Error().Msg("[CollectHandler][Create] MustGet(uid) error")
		return
	}

	tid, ok := c.MustGet("tid").(string)
	if !ok {
		logger.Error().Msg("[CollectHandler][Create] MustGet(tid) error")
		return
	}

	_, err := uuid.FromString(tid)
	if err != nil {
		logger.Error().Any("tid", tid).Msg("tid is not uuid")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	_, err = uuid.FromString(uid)
	if err != nil {
		logger.Error().Any("uid", uid).Msg("uid is not uuid")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	collect, err := ch.collectStore.GetCollectByUidTid(c, uid, tid)
	if err != nil {
		logger.Error().Any("collect", collect).Msg("collect not exist")
		handler.SendResponse(c, errmsg.ERROR_CREATE_ALREADY_COLLECT, nil)
		return
	}

	err = ch.collectStore.DeleteCollect(c, tid, uid)
	if err != nil {
		logger.Error().Err(err).Msgf("tid: %s, uid: %s", tid, uid)
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	handler.SendResponse(c, errmsg.SUCCSE, nil)
}
