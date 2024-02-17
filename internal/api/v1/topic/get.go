package topic

import (
	"time"
)

type GetTopicResp struct {
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
	RefundEndTime time.Time `json:"refund_end_time"`
	EndTime       time.Time `json:"end_time"`
	// Category      *core.Category `json:"category"`
}

// func (th *TopicHandler) Get(c *gin.Context) {
// 	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

// 	tid, err := th.checkTid(c)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("[Get][checkCreate]")
// 		handler.SendResponse(c, errmsg.ERROR, nil)
// 		return
// 	}

// 	getTopicResp, err := th.getTopicResp(c, tid)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("[Get][storage.GetTopicByTid]")
// 		handler.SendResponse(c, errmsg.ERROR, nil)
// 		return
// 	}
// 	logger.Info().Any("getTopicResp", getTopicResp).Msg("[Get]")

// 	handler.SendResponse(c, errmsg.SUCCSE, getTopicResp)
// }

// func (th *TopicHandler) getTopicResp(c *gin.Context, tid string) (*GetTopicResp, error) {
// 	topic, err := th.topicStore.GetTopicByTid(c, tid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	getTopicResp := new(GetTopicResp)
// 	copier.Copy(&getTopicResp, &topic)
// 	getTopicResp.Category = th.categoryMap[getTopicResp.Cid]
// 	return getTopicResp, nil
// }
