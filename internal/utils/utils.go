package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
)

func GetOrGenUid(c *gin.Context) string {
	xid, ok := c.Get(consts.DefaultXid)
	if !ok {
		xid = convert.NewUUID()
	}
	return xid.(string)
}

func GetUserId(c *gin.Context) (string, bool) {
	uidS, ok := c.Get("uid")
	if !ok {
		return "", false
	}

	uid, ok := uidS.(string)
	if !ok {
		return "", false
	}
	return uid, true
}
