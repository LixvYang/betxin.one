package category

import (
	"github.com/gin-gonic/gin"
)

type CreateCategoryReq struct {
	Name string `json:"name"`
}

func (ch *CategoryHandler) Create(c *gin.Context) {
	// logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	// var req CreateCategoryReq

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	logger.Error().Err(err).Msg("[Create][ShouldBindJSON] error")
	// 	handler.SendResponse(c, errmsg.ERROR_BIND, nil)
	// 	return
	// }

	// err := ch.storage.CreateCategory(c, logger, req.Name)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("[Create][CreateCategory] error")
	// 	handler.SendResponse(c, errmsg.ERROR, nil)
	// 	return
	// }
	// handler.SendResponse(c, errmsg.SUCCSE, nil)
}
