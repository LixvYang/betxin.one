package topic

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/internal/utils/token"

	"github.com/rs/zerolog"
)

type ListTopicByCidReq struct {
	Uid       string `json:"uid"`
	PageToken string `json:"page_token"`
}

type ListTopicByCidResp struct {
	List         []*ListTopicData `json:"list"`
	PrePageToken string           `json:"pre_page_token"`
}

type ListTopicData struct {
	*schema.Topic
	IsCollect bool             `json:"is_collect"`
	Category  *schema.Category `json:"category"`
}

const (
	defaultPageSize int64 = 10
	defaultCursor   int64 = 10
)

func (t *TopicHandler) ListTopicsByCid(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	cidS := c.Query("cid")
	if cidS == "" {
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	cid, err := strconv.ParseInt(cidS, 0, 64)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	var req ListTopicByCidReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	page := token.Token(req.PageToken).Decode()
	var (
		cursor   int64 = defaultCursor
		pageSize int64 = defaultPageSize
	)

	if req.PageToken != "" {
		// 解析分页
		if page.NextTimeAtUTC > time.Now().Unix() || time.Now().Unix()-page.NextTimeAtUTC > int64(time.Hour)*24 {
			logger.Error().Msgf("bad page token, key: %s", req.PageToken)
			handler.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		// invaild
		if page.PreID <= 0 || page.NextTimeAtUTC == 0 || page.NextTimeAtUTC > time.Now().Unix() || page.PageSize <= 0 {
			logger.Error().Msgf("bad page token, key: %s", req.PageToken)
			return
		}
		cursor = page.PreID
		pageSize = page.PageSize
	}

	topicList, err := t.storage.ListTopicByCid(c, logger, cid, cursor, pageSize+1)
	if err != nil {
		logger.Error().Err(err).Msgf("[t.ListTopicsByCid][ListTopicByCid] err")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	topicDataList := t.getTopicDataList(c, topicList)

	var (
		hasPrePage   bool
		prePageToken string
	)

	if len(topicDataList) > int(pageSize) {
		hasPrePage = true
	}

	// if has pre page
	if hasPrePage {
		prePageInfo := token.Page{
			PreID:         topicDataList[len(topicDataList)-1].Tid,
			NextTimeAtUTC: time.Now().Unix(),
			PageSize:      pageSize,
		}
		prePageToken = string(prePageInfo.Encode())

		handler.SendResponse(c, errmsg.SUCCSE, &ListTopicByCidResp{
			PrePageToken: prePageToken,
			List:         topicDataList[:len(topicDataList)-1],
		})
		return
	}

	handler.SendResponse(c, errmsg.SUCCSE, &ListTopicByCidResp{
		PrePageToken: prePageToken,
		List:         topicDataList,
	})
}

func (t *TopicHandler) getTopicDataList(c *gin.Context, args []*schema.Topic) []*ListTopicData {
	topicDataList := make([]*ListTopicData, len(args))
	copier.Copy(&topicDataList, &args)
	for i := 0; i < len(topicDataList); i++ {
		topicDataList[i].Category = t.categoryMap[topicDataList[i].Cid]
	}
	uidS, exists := c.Get("uid")
	if !exists {
		return topicDataList
	}

	uid, ok := uidS.(string)
	if !ok {
		return topicDataList
	}

	collects, err := t.topicCollectHandler.GetCollectByUid(uid)
	if err != nil {
		return topicDataList
	}

	collectMap := make(map[int64]bool)
	for _, collect := range collects {
		collectMap[collect.Tid] = collect.Status
	}

	for i := 0; i < len(topicDataList); i++ {
		topicDataList[i].IsCollect = collectMap[topicDataList[i].Tid]
	}

	return topicDataList
}
