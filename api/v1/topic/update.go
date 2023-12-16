package topic

import (
	"github.com/gin-gonic/gin"
)

type UpdateTopicInfoReq struct {
	Cid           int64  `json:"cid"`
	Title         string `json:"title"`
	Intro         string `json:"intro"`
	Content       string `json:"content"`
	ImgURL        string `json:"img_url"`
	RefundEndTime int64  `json:"refund_end_time"`
	EndTime       int64  `json:"end_time"`
}

func (th *TopicHandler) UpdateTopicInfo(c *gin.Context) {
	// logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)

	// var req UpdateTopicInfoReq
	// err := c.ShouldBindJSON(&req)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("[UpdateTopicInfo][ShouldBindJSON] error")
	// 	handler.SendResponse(c, errmsg.ERROR_BIND, nil)
	// 	return
	// }

	// err = th.checkUpdateTopicInfoReq(c, &req)
	// if err != nil {
	// 	logger.Error().Any("req", req).Str("tid", c.Param("tid")).Err(err).Msg("[UpdateTopicInfo][checkUpdateTopicInfoReq] err")
	// 	handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
	// 	return
	// }
	// schemaTopic := new(schema.Topic)
	// copier.Copy(schemaTopic, &req)
	// err = th.storage.UpdateTopicInfo(c, logger, schemaTopic)
	// if err != nil {
	// 	logger.Error().Any("req", req).Str("tid", c.Param("tid")).Err(err).Msg("[UpdateTopicInfo][storage.UpdateTopicInfo] err")
	// 	handler.SendResponse(c, errmsg.ERROR_INVAILD_ARGV, nil)
	// 	return
	// }
	// handler.SendResponse(c, errmsg.SUCCSE, req)
}

// func (th *TopicHandler) checkUpdateTopicInfoReq(c *gin.Context, req *UpdateTopicInfoReq) error {
// 	tid, err := convert.StrToInt64(c.Param("tid"))
// 	if err != nil || tid == 0 {
// 		return errors.New("[checkUpdateTopicInfoReq][StrToInt64] tid invalid")
// 	}

// 	createTopicReq := new(CreateTopicReq)
// 	copier.Copy(createTopicReq, req)
// 	return th.checkCreateReq(createTopicReq)
// }
