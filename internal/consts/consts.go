package consts

import "time"

const (
	DefaultLoggerKey = "logger"
	DefaultXid       = "X-Reqid"

	Uid = "uid"
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

// mongo coll name
const (
	UserCollection          = "user"
	CategoryCollection      = "category"
	RefundCollection        = "refund"
	TopicCollection         = "topic"
	TopicPurchaseCollection = "topic_purchase"
	BonuseCollection        = "bonuse"
)

const (
	// 延时双删除时间
	DelayedDeletionInterval = time.Second >> 1
)
