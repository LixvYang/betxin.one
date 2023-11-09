package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/pkg/jwt"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild token",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild token",
			})
			c.Abort()
			return
		}

		mc, err := jwt.ParseJwt(parts[1])
		if err != nil {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "invaild token",
			})
			c.Abort()
			return
		}
		c.Set("uid", mc.Uid)
		c.Next()
	}
}
