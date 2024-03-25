package topicpurchase

import (
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type QueryPurchaseReq struct {
	RequestID string `uri:"request_id" binding:"required"`
}

func (tph *TopicPurchaseHandler) QueryTopicPurchaseHistory(ctx *gin.Context) {
	logger := ctx.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	var req QueryPurchaseReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		logger.Error().Err(err).Msg("invalid request")
		handler.SendResponse(ctx, errmsg.ERROR, nil)
		return
	}

	_, err := convert.VaildUUID(req.RequestID)
	if err != nil {
		logger.Error().Err(err).Msg("invalid request_id")
		handler.SendResponse(ctx, errmsg.ERROR, nil)
		return
	}

	_, err = tph.storage.QueryTopicPurchaseHistory(ctx, req.RequestID)
	if err != nil {
		if err == mongo.ErrNoSuchItem {
			logger.Error().Any("request_id", req.RequestID).Err(err).Msg("query topic purchase history failed")
			handler.SendResponse(ctx, errmsg.ERROR_BUY_RECORD_EMPTY, nil)
			return
		}
		logger.Error().Any("request_id", req.RequestID).Err(err).Msg("query topic purchase history failed")
		handler.SendResponse(ctx, errmsg.ERROR, nil)
		return
	}
	handler.SendResponse(ctx, errmsg.SUCCES, nil)
}
