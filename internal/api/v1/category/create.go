package category

import (
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CreateCategoryReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (ch *CategoryHandler) Create(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	var req CreateCategoryReq

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().Err(err).Msg("[Create][ShouldBindJSON] error")
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	err := ch.storage.CreateCategory(c, &schema.Category{
		ID:   req.ID,
		Name: req.Name,
	})

	if err != nil {
		logger.Error().Err(err).Msg("[Create][CreateCategory] error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	handler.SendResponse(c, errmsg.SUCCES, nil)
}
