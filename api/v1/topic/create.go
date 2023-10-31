package topic

import "github.com/gin-gonic/gin"

func (t *TopicHandler) Create(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "world",
	})
}
