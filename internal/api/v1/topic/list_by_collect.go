package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type ListTopicsByUidCollectResp struct {
	List      []*ListTopicData `json:"list"`
	Total     int64            `json:"total"`
	ConnectId string           `json:"connect_id"`
}

type ListTopicsByUidCollectReq struct {
	Uid    string `json:"-"`
	Limit  int64  `form:"limit"`
	Offset int64  `form:"offset"`
}

func checkListCollectsByUidReq(c *gin.Context) (*ListTopicsByUidCollectReq, error) {
	var req ListTopicsByUidCollectReq
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	uidS, exists := c.Get(consts.Uid)
	if !exists {
		return nil, consts.ErrUidNotExist
	}

	uid := uidS.(string)

	if req.Limit < 0 {
		req.Limit = consts.DefaultLimit
	}

	if req.Offset < 0 {
		req.Offset = consts.DefaultOffset
	}

	req.Uid = uid
	return &req, nil
}

// 查询用户收藏的话题
func (th *TopicHandler) ListTopicsByCollect(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	req, err := checkListCollectsByUidReq(c)
	if err != nil {
		logger.Error().Err(err).Msg("[CollectHandler][ListTopicsByCollect][checkListCollectsByUidReq]")
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	collects, total, err := th.collectSrv.GetCollectsByUid(c, &logger, req.Uid, req.Limit, req.Offset)
	if err != nil {
		logger.Error().Err(err).Msg("[th.ListTopicsByCollect][GetCollectByUid]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	tids := make([]string, len(collects))
	lo.ForEach(collects, func(item *schema.Collect, i int) {
		tids[i] = item.Tid
	})

	topics, err := th.topicSrv.GetTopicsByTids(c, &logger, tids)
	if err != nil {
		logger.Error().Err(err).Msg("[th.ListTopicsByCollect][GetTopicsByTids]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	respTopics := th.getTopicDataList(c, &logger, topics)

	var resp ListTopicsByUidCollectResp
	copier.Copy(&resp.List, &respTopics)
	resp.Total = total
	resp.ConnectId = c.GetString(consts.DefaultXid)

	handler.SendResponse(c, errmsg.SUCCSE, resp)
}
