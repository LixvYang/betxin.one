package topic

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql/core"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/lixvyang/betxin.one/internal/utils/token"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ListTopicByCidResp struct {
	List         []*ListTopicData `json:"list"`
	PrePageToken string           `json:"pre_page_token"`
}

type TopicItem struct {
	Tid           string    `json:"tid"`
	Cid           int64     `json:"cid"`
	Title         string    `json:"title"`
	Intro         string    `json:"intro"`
	Content       string    `json:"content"`
	YesRatio      string    `json:"yes_ratio"`
	NoRatio       string    `json:"no_ratio"`
	YesCount      string    `json:"yes_count"`
	NoCount       string    `json:"no_count"`
	TotalCount    string    `json:"total_count"`
	CollectCount  int64     `json:"collect_count"`
	ReadCount     int64     `json:"read_count"`
	ImgURL        string    `json:"img_url"`
	IsStop        bool      `json:"is_stop"`
	IsDeleted     bool      `json:"is_deleted"`
	RefundEndTime int64     `json:"refund_end_time"`
	EndTime       int64     `json:"end_time"`
	CreatedAt     time.Time `json:"created_at"`
}

type ListTopicData struct {
	*TopicItem
	Tid       string         `json:"tid"`
	IsCollect bool           `json:"is_collect"`
	Category  *core.Category `json:"category"`
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
		cursor   time.Time = time.Now()
		pageSize int64     = defaultPageSize
	)

	if pageToken != "" {
		// 解析分页
		if page.NextTimeAtUTC < time.Now().UnixMilli() || time.Now().Unix()-page.NextTimeAtUTC > int64(time.Hour)*24 {
			logger.Error().Msgf("bad page token invaild page time: %#v", page)
			handler.SendResponse(c, errmsg.ERROR, nil)
			return
		}

		// invaild
		if page.PageSize <= 0 {
			logger.Error().Msgf("bad page token invaild page info, page: %#v", page)
			return
		}
		cursor = page.CreatedAt
		pageSize = page.PageSize
	}

	topicList, err := t.topicStore.ListTopicsByCid(c, cid, cursor, pageSize+1)
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
			CreatedAt:     topicDataList[len(topicDataList)-1].CreatedAt,
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

func (t *TopicHandler) getTopicDataList(c *gin.Context, args []*core.Topic) []*ListTopicData {
	topicDataList := make([]*ListTopicData, len(args))
	copier.Copy(&topicDataList, &args)
	for i := 0; i < len(topicDataList); i++ {
		topicDataList[i].Tid = args[i].Tid
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

	collects, err := t.collectStore.GetCollectByUid(c, uid)
	if err != nil {
		return topicDataList
	}

	collectMap := make(map[string]bool)
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
		return 0, errors.New("cid error")
	}

	cid, err := strconv.ParseInt(cidS, 0, 64)
	if err != nil {
		return 0, errors.New("cid convert error")
	}

	_, ok := t.categoryMap[cid]
	if !ok {
		return 0, errors.New("cid not exist")
	}
	return cid, nil
}
