package category

import (
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (ch *CategoryHandler) Delete(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	id, err := ch.checkCategoryId(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][checkCategoryId]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = ch.storage.DeleteCategory(c, &logger, id)
	if err != nil {
		logger.Error().Err(err).Msg("[Delete][storage.GetCategoryById]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Int64("tid", id).Msg("[Delete]")
	handler.SendResponse(c, errmsg.SUCCSE, nil)
}

func (ch *CategoryHandler) checkCategoryId(c *gin.Context) (int64, error) {
	id, err := convert.StrToInt64(c.Param("id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}
