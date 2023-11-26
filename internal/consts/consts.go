package consts

import "time"

const (
	LoggerKey = "logger"
	Xid       = "xid"
	Uid       = "uid"
)

const (
	CacheNotFound     = "cache no data"
	CacheAlreadyExist = "data already exist"
)

// redis前缀
const (
	// hset betxin_user_info uid {{ info }}
	RdsHashUserInfoKey = "betxin_user_info"
	// hset betxin_topic_info tid {{ info }}
	RdsHashTopicInfoKey = "betxin_topic_info"
)

const (
	// 延时双删除时间
	DelayedDeletionInterval = time.Second >> 1
)
