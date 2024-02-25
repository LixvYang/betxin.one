package consts

import (
	"errors"
	"time"
)

const (
	DefaultLoggerKey = "logger"
	DefaultXid       = "X-Reqid"

	Uid = "uid"

	DefaultLimit  int64 = 10
	DefaultOffset int64 = 0
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
	UserCollection                 = "user"
	CategoryCollection             = "category"
	CollectCollection              = "collect"
	RefundCollection               = "refund"
	TopicCollection                = "topic"
	TopicPurchaseCollection        = "topic_purchase"
	TopicPurchaseHistoryCollection = "topic_purchase_history"
	BonuseCollection               = "bonuse"
	MixinUtxoCollection            = "mixin_utxo"
)

const (
	// 延时双删除时间
	DelayedDeletionInterval = time.Second >> 1
)

var (
	ErrUidNotExist = errors.New("uid not found")
)
