package user

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lixvyang/betxin.one/api/v1"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
)

func (u *UserHandler) Create(c *gin.Context) {
	v1.SendResponse(c, errmsg.SUCCSE, gin.H{
		"Hello": "world",
	})
	return
}
