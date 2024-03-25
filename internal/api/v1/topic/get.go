package topic

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/internal/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/model/database/schema"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

type GetTopicResp struct {
	Tid           string           `json:"tid"`
	Cid           int64            `json:"cid"`
	Title         string           `json:"title"`
	Intro         string           `json:"intro"`
	Content       string           `json:"content"`
	YesRatio      string           `json:"yes_ratio"`
	NoRatio       string           `json:"no_ratio"`
	YesCount      string           `json:"yes_count"`
	NoCount       string           `json:"no_count"`
	TotalCount    string           `json:"total_count"`
	CollectCount  int64            `json:"collect_count"`
	ReadCount     int64            `json:"read_count"`
	ImgURL        string           `json:"img_url"`
	RefundEndTime time.Time        `json:"refund_end_time"`
	EndTime       time.Time        `json:"end_time"`
	Category      *schema.Category `json:"category"`
	IsCollect     bool             `json:"is_collect"`
}

func (th *TopicHandler) GetTopicByTid(c *gin.Context) {
	logger := c.MustGet(consts.DefaultLoggerKey).(zerolog.Logger)

	tid, err := th.checkTid(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][checkCreate]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	getTopicResp, err := th.getTopicResp(c, &logger, tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][storage.GetTopicByTid]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	logger.Info().Any("getTopicResp", getTopicResp).Msg("[Get]")

	handler.SendResponse(c, errmsg.SUCCES, getTopicResp)
}

func (th *TopicHandler) getTopicResp(c *gin.Context, logger *zerolog.Logger, tid string) (*GetTopicResp, error) {
	topic, err := th.storage.GetTopicByTid(c, tid)
	if err != nil {
		return nil, err
	}
	getTopicResp := new(GetTopicResp)
	copier.Copy(&getTopicResp, &topic)
	category, err := th.storage.GetCategoryById(c, topic.Cid)
	if err != nil {
		logger.Error().Err(err).Msg("[getTopicResp][storage.GetCategoryById]")
		return nil, err
	}
	getTopicResp.Category = category
	getTopicResp.TotalCount = convert.DecimalAdd(getTopicResp.YesCount, getTopicResp.NoCount).String()
	getTopicResp.YesRatio = convert.DecimalDiv(getTopicResp.YesCount, getTopicResp.TotalCount).String()
	getTopicResp.NoRatio = convert.DecimalDiv(getTopicResp.NoCount, getTopicResp.TotalCount).String()
	if getTopicResp.YesCount == getTopicResp.NoCount {
		getTopicResp.YesRatio = "50"
		getTopicResp.NoRatio = "50"
	}

	uidS, exists := c.Get("uid")
	if !exists {
		return getTopicResp, nil
	}

	uid, ok := uidS.(string)
	if !ok {
		return getTopicResp, nil
	}

	collects, err := th.storage.ListCollects(c, uid)
	if err != nil {
		return getTopicResp, nil
	}

	collectMap := make(map[string]bool)
	for _, collect := range collects {
		collectMap[collect.Tid] = collect.Status
	}

	getTopicResp.IsCollect = collectMap[tid]
	return getTopicResp, nil
}
