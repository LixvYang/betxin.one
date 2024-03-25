package topic

import (
	"errors"

	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mongo"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ListTopicsReq struct {
	C      int64 `form:"_c"`
	Limit  int64 `form:"limit"`
	Offset int64 `form:"offset"`

	cid int64
}

type ListTopicsResp struct {
	List  []*ListTopicData `json:"list"`
	Total int64            `json:"total"`
}

func (th *TopicHandler) ListTopics(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)
	req, err := th.checkListTopicsReq(c)
	if err != nil {
		logger.Error().Err(err).Msg("check list topics req error")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	topics, total, err := th.storage.ListTopics(c, req.cid, req.Limit, req.Offset)
	if err != nil {
		logger.Error().Err(err).Msg("list topics error")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	topicListResp := th.getTopicDataList(c, &logger, topics)

	handler.SendResponse(c, errmsg.SUCCES, &ListTopicsResp{
		List:  topicListResp,
		Total: total,
	})
}

func (th *TopicHandler) checkListTopicsReq(c *gin.Context) (*ListTopicsReq, error) {
	var req ListTopicsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if req.C != 0 {
		// 验证category
		_, err := th.storage.GetCategoryById(c, req.C)
		if err != nil {
			if err == mongo.ErrNoSuchItem {
				return nil, errors.New("category not found")
			}
			return nil, err
		}
		req.cid = req.C
	}

	const (
		defaultLimit  int64 = 10
		defaultOffset int64 = 0
	)

	if req.Limit <= 0 {
		req.Limit = defaultLimit
	}

	if req.Offset < 0 {
		req.Offset = defaultOffset
	}

	return &req, nil
}
