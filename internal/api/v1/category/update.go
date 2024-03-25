package category

import (
	"errors"

	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UpdateCategoryReq struct {
	Name string `json:"name"`
}

func (ch *CategoryHandler) Update(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(*zerolog.Logger)

	var req UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Err(err).Any("req", req).Msg("[Update][ShouldBindJSON] err")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	id, err := ch.checkUpdateReq(c, logger, req.Name)
	if err != nil {
		logger.Error().Err(err).Msg("[Update][storage.checkUpdateReq]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	err = ch.storage.UpdateCategory(c, id, req.Name)
	if err != nil {
		logger.Error().Err(err).Msg("[Update][storage.GetCategoryById]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Any("name", req.Name).Msg("[Update] success")
	handler.SendResponse(c, errmsg.SUCCES, nil)
}

func (ch *CategoryHandler) checkUpdateReq(c *gin.Context, logger *zerolog.Logger, name string) (int64, error) {
	id, err := ch.checkCategoryId(c)
	if err != nil {
		return 0, errors.New("checkCategoryId error")
	}

	_, err = ch.storage.GetCategoryById(c, id)
	if err != nil {
		return 0, err
	}

	if name == "" {
		return 0, errors.New("name invalied")
	}
	return id, nil
}
