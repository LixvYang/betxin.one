package category

import (
	"github.com/gin-gonic/gin"
)

func (ch *CategoryHandler) Get(c *gin.Context) {
	// logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	// id, err := ch.checkCategoryId(c)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("[Get][checkCategoryId]")
	// 	handler.SendResponse(c, errmsg.ERROR, nil)
	// 	return
	// }

	// category, err := ch.storage.GetCategoryById(c, logger, id)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("[Get][storage.GetCategoryById]")
	// 	handler.SendResponse(c, errmsg.ERROR, nil)
	// 	return
	// }
	// logger.Info().Any("category", category).Msg("[Get]")
	// handler.SendResponse(c, errmsg.SUCCSE, category)
}
