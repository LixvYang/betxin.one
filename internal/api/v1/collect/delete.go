package collect

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

type DeleteCollectReq struct {
	Tid string `json:"tid"`
	uid string `json:"-"`
}

func checkDeleteCollectReq(c *gin.Context) (*DeleteCollectReq, error) {
	var req DeleteCollectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	req.uid = c.GetString(consts.Uid)
	return &req, nil
}

func (ch *CollectHandler) Delete(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	req, err := checkDeleteCollectReq(c)
	if err != nil {
		logger.Error().Err(err).Msg("checkDeleteCollectReq error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	// 判断tid是否存在
	_, err = ch.topicSrv.GetTopicByTid(c, &logger, req.Tid)
	if err != nil {
		logger.Error().Err(err).Msg("get topic by tid err")
		handler.SendResponse(c, errmsg.ERROR_GET_TOPIC, nil)
		return
	}

	collect, err := ch.storage.GetCollectByUidTid(c, &logger, req.uid, req.Tid)
	if err != nil {
		if err == mongo.ErrNoSuchCollect {
			logger.Error().Err(err).Msg("get collect by uid tid err")
			handler.SendResponse(c, errmsg.ERROR_NOT_COLLECT, nil)
			return
		}
		logger.Error().Err(err).Msg("get collect by uid tid err")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	collect.Status = false
	collect.UpdatedAt = time.Now()

	if err := ch.storage.UpsertCollect(c, &logger, req.uid, req.Tid, collect); err != nil {
		logger.Error().Any("req", req).Err(err).Msg("DeleteCollect error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	logger.Info().Any("req", req).Msg("DeleteCollect success")
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}
