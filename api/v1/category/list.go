package category

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

func (ch *CategoryHandler) List(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	categories, err := ch.storage.ListCategories()
	if err != nil {
		logger.Error().Err(err).Msg("[List][ListCategories]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Any("categories", categories).Msg("[List][ListCategories]")

	handler.SendResponse(c, errmsg.SUCCSE, categories)
}
