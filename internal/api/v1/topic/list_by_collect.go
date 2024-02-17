package topic

// 查询用户收藏的话题
// func (th *TopicHandler) ListTopicsByCollect(c *gin.Context) {
// 	logger := c.MustGet(consts.LoggerKey).(*zerolog.Logger)
// 	uid, ok := c.MustGet("uid").(string)
// 	if !ok {
// 		logger.Error().Msg("[CollectHandler][Create] MustGet(uid) error")
// 		return
// 	}
// 	_, err := uuid.FromString(uid)
// 	if err != nil {
// 		logger.Error().Any("uid", uid).Msg("uid is not uuid")
// 		handler.SendResponse(c, errmsg.ERROR, nil)
// 		return
// 	}

// 	collects, err := th.collectStore.GetCollectByUid(c, uid)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("[th.ListTopicsByCollect][GetCollectByUid]")
// 		handler.SendResponse(c, errmsg.ERROR, nil)
// 		return
// 	}

// 	tids := make([]string, len(collects))
// 	lo.ForEach(collects, func(item *core.Collect, i int) {
// 		tids[i] = item.Tid
// 	})

// 	topics, err := th.topicStore.GetTopicsByTids(c, tids)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("[th.ListTopicsByCollect][GetTopicsByTids]")
// 		handler.SendResponse(c, errmsg.ERROR, nil)
// 		return
// 	}

// 	resp := th.getTopicDataList(c, topics)
// 	handler.SendResponse(c, errmsg.SUCCSE, resp)
// }
