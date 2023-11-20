package refund

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/rs/zerolog"
)

func Create(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	uid := c.MustGet("uid").(string)
	if uid == "" {
		logger.Error().Send()
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	/*
		1. 退款
		2. 获取uid
		3. 判断话题停止退款时间
		4. 获取uid的购买记录 判断请求退款数是否对的上，判断话题时间
		5. 向用户转账 并且发送消息
		
	*/

}
