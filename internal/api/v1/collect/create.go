package collect

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

type CreateCollectParams struct {
	Tid string `json:"tid" binding:"required"`
	Uid string `json:"-"`
}

func (ch *CollectHandler) Create(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	req, err := checkCreateCollect(c)
	if err != nil {
		logger.Error().Err(err).Msg("check params err")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, err.Error())
		return
	}

	// 判断tid是否存在
	_, err = ch.storage.GetTopicByTid(c, req.Tid)
	if err != nil {
		logger.Error().Any("req", req).Err(err).Msg("get topic by tid err")
		handler.SendResponse(c, errmsg.ERROR_GET_TOPIC, nil)
		return
	}

	argv := schema.Collect{
		Tid:       req.Tid,
		UID:       req.Uid,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ch.storage.UpsertCollect(c, req.Uid, argv.Tid, &argv)
	if err != nil {
		if err == mongo.ErrItemExist {
			logger.Error().Err(err).Msg("collect already exists")
			handler.SendResponse(c, errmsg.ERROR_CREATE_ALREADY_COLLECT, nil)
			return
		}
		logger.Error().Any("argv", argv).Err(err).Msg("create collect err")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	handler.SendResponse(c, errmsg.SUCCES, nil)
}

func checkCreateCollect(c *gin.Context) (*CreateCollectParams, error) {
	uidS, exists := c.Get("uid")
	if !exists {
		return nil, errors.New("uid not exits")
	}

	uid := uidS.(string)

	var req CreateCollectParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}
	if req.Tid == "" {
		return nil, errors.New("tid is empty")
	}
	req.Uid = uid

	return &req, nil
}
