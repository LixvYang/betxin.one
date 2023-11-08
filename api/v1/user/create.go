package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
)

func (u *UserHandler) Create(c *gin.Context) {
	handler.SendResponse(c, errmsg.SUCCSE, gin.H{
		"Hello": "world",
	})
}
