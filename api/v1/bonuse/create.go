package bonuse

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (bh *BonuseHandler) Create(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
	uid := c.MustGet("uid").(string)
	if uid == "" {
		logger.Error().Send()
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	/*
		1. 创建奖金
		2. 获取tid
		3. 拿到话题购买列表 表中的赢了的人的uid
		4. 给每个人发送奖金
	*/

}
