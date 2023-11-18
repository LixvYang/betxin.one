package topic

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/internal/utils/token"
	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

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
	defaultCursor   int64 = math.MaxInt64
)

func (t *TopicHandler) ListTopicsByCid(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	cid, err := t.checkListTopicsByCidReq(c)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
		return
	}

	pageToken := c.Query("page_token")
	page := token.Token(pageToken).Decode()
	var (
		cursor   int64 = defaultCursor
		pageSize int64 = defaultPageSize
	)

	if pageToken != "" {
		// 解析分页
		if page.NextTimeAtUTC < time.Now().UnixMilli() || time.Now().Unix()-page.NextTimeAtUTC > int64(time.Hour)*24 {
			logger.Error().Msgf("bad page token invaild page time: %#v", page)
			handler.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		// invaild
		if page.PreID <= 0 || page.PageSize <= 0 {
			logger.Error().Msgf("bad page token invaild page info, page: %#v", page)
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
			NextTimeAtUTC: time.Now().Add(time.Hour * 24).UnixMilli(),
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

func (t *TopicHandler) checkListTopicsByCidReq(c *gin.Context) (int64, error) {
	cidS := c.Param("cid")
	if cidS == "" {
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return 0, errors.New("cid error")
	}

	cid, err := strconv.ParseInt(cidS, 0, 64)
	if err != nil {
		handler.SendResponse(c, errmsg.ERROR_BIND, nil)
		return 0, errors.New("cid convert error")
	}

	_, ok := t.categoryMap[cid]
	if !ok {
		handler.SendResponse(c, errmsg.ERROR, nil)
		return 0, errors.New("cid not exist")
	}
	return cid, nil
}
