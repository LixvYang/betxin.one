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

/*
1. 传递uid/不传递uid都可以 要求使用同样的api 传递 (要求中间件那里不做限制)
2. 传递 uid有更多的可选动作
*/
func JWTAuthNotMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
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
		}
		c.Next()
	}
}
