package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/lixvyang/betxin.one/api/v1"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/pkg/jwt"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			v1.SendResponse(c, errmsg.ERROR_AUTH, nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			v1.SendResponse(c, errmsg.ERROR_INVAILD_TOKEN, nil)
			c.Abort()
			return
		}

		mc, err := jwt.ParseJwt(parts[1])
		if err != nil {
			v1.SendResponse(c, errmsg.ERROR_INVAILD_TOKEN, nil)
			c.Abort()
			return
		}
		c.Set("uid", mc.Uid)
		c.Next()
	}
}
