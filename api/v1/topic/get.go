package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/api/v1/handler"
	"github.com/lixvyang/betxin.one/internal/consts"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

type GetTopicResp struct {
	Tid           int64  `json:"tid"`
	Cid           int64  `json:"cid"`
	Title         string `json:"title"`
	Intro         string `json:"intro"`
	Content       string `json:"content"`
	YesRatio      string `json:"yes_ratio"`
	NoRatio       string `json:"no_ratio"`
	YesCount      string `json:"yes_count"`
	NoCount       string `json:"no_count"`
	TotalCount    string `json:"total_count"`
	CollectCount  int64  `json:"collect_count"`
	ReadCount     int64  `json:"read_count"`
	ImgURL        string `json:"img_url"`
	IsStop        bool   `json:"is_stop"`
	IsDeleted     bool   `json:"is_deleted"`
	RefundEndTime int64  `json:"refund_end_time"`
	EndTime       int64  `json:"end_time"`
}

func (th *TopicHandler) Get(c *gin.Context) {
	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	tid, err := th.checkTid(c)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][checkCreate]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	topic, err := th.storage.GetTopicByTid(c, logger, tid)
	if err != nil {
		logger.Error().Err(err).Msg("[Get][storage.GetTopicByTid]")
		handler.SendResponse(c, errmsg.ERROR, nil)
		return
	}
	getTopicResp := new(GetTopicResp)
	copier.Copy(getTopicResp, topic)
	logger.Info().Any("topic", topic).Msg("[Get]")
	
	handler.SendResponse(c, errmsg.SUCCSE, getTopicResp)
}
