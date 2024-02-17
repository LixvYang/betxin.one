package topicpurchase

import (
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateReq struct {
	Tid string
}

// 话题购买
func (tph *TopicPurchaseHandler) Create(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	uid, ok := c.MustGet("uid").(string)
	if !ok || uid == "" {
		logger.Error().Msg("uid not exist")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	// 1. 获取uid

}
