package topic

import (
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ListTopicsReq struct {
	Category string `json:"_c"`
	Limit    int64  `form:"limit"`
	Offset   int64  `form:"offset"`

	cid int64
}

type ListTopicsResp struct {
	List      []*schema.Topic `json:"list"`
	Total     int64           `json:"total"`
	ConnectId string          `json:"connect_id"`
}

func (th *TopicHandler) ListTopics(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	connId := c.GetString(consts.DefaultXid)

	req, err := th.checkListTopicsReq(c, &logger)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	topics, total, err := th.topicSrv.ListTopics(c, &logger, req.cid, req.Limit, req.Offset)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	handler.SendResponse(c, errmsg.SUCCSE, &ListTopicsResp{
		List:      topics,
		Total:     total,
		ConnectId: connId,
	})
}

func (th *TopicHandler) checkListTopicsReq(c *gin.Context, logger *zerolog.Logger) (*ListTopicsReq, error) {
	var req ListTopicsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	// 验证category
	// category, err := th.categorySrv.ListCategories(c, logger)
	// if err != nil {
	// 	return err
	// }
	req.cid = 0

	return &req, nil
}
