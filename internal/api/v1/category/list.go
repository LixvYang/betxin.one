package category

import (
	"github.com/gin-gonic/gin"
)

func (ch *CategoryHandler) List(c *gin.Context) {
	// logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	// categories, err := ch.storage.ListCategories()
	// if err != nil {
	// 	logger.Error().Err(err).Msg("[List][ListCategories]")
	// 	handler.SendResponse(c, errmsg.ERROR, nil)
	// 	return
	// }
	// logger.Info().Any("categories", categories).Msg("[List][ListCategories]")

	// handler.SendResponse(c, errmsg.SUCCSE, categories)
}
