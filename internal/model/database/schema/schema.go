package schema

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	IdentityNumber string    `bson:"identity_number" json:"identity_number"`
	Uid            string    `bson:"uid" json:"uid"`
	FullName       string    `bson:"full_name" json:"full_name"`
	AvatarURL      string    `bson:"avatar_url" json:"avatar_url"`
	SessionID      string    `bson:"session_id" json:"session_id"`
	Biography      string    `bson:"biography" json:"biography"`
	PrivateKey     string    `bson:"private_key" json:"private_key"`
	ClientID       string    `bson:"client_id" json:"client_id"`
	Contract       string    `bson:"contract" json:"contract"`
	IsMvmUser      bool      `bson:"is_mvm_user" json:"is_mvm_user"`
	MvmPublicKey   string    `bson:"mvm_public_key" json:"mvm_public_key"`
	MixinCreatedAt time.Time `bson:"mixin_created_at" json:"mixin_created_at"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

type Topic struct {
	Tid           string    `bson:"tid" json:"tid"`
	Cid           int64     `bson:"cid" json:"cid"`
	Title         string    `bson:"title" json:"title"`
	Intro         string    `bson:"intro" json:"intro"`
	Content       string    `bson:"content" json:"content"`
	YesAmount     string    `bson:"yes_amount" json:"yes_amount"`
	NoAmount      string    `bson:"no_amount" json:"no_amount"`
	CollectCount  int64     `bson:"collect_count" json:"collect_count"`
	ReadCount     int64     `bson:"read_count" json:"read_count"`
	ImgURL        string    `bson:"img_url" json:"img_url"`
	IsStop        bool      `bson:"is_stop" json:"is_stop"`
	RefundEndTime time.Time `bson:"refund_end_time" json:"refund_end_time"`
	EndTime       time.Time `bson:"end_time" json:"end_time"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt     time.Time `bson:"deleted_at" json:"deleted_at"`
}

type Category struct {
	ID   int64  `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type Bonuse struct {
	Uid       string    `bson:"uid" json:"uid"`
	Tid       string    `bson:"tid" json:"tid"`
	AssetId   string    `bson:"asset_id" json:"asset_id"`
	Amount    string    `bson:"amount" json:"amount"`
	Memo      string    `bson:"memo" json:"memo"`
	TraceId   string    `bson:"trace_id" json:"trace_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at" json:"deleted_at"`
}

type Collect struct {
	UID       string    `json:"uid" bson:"uid"`
	Tid       string    `json:"tid" bson:"tid"`
	Status    bool      `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// type Feedback struct {
// 	UID       string `json:"uid"`
// 	Fid       string `json:"fid"`
// 	Title     string `json:"title"`
// 	Content   string `json:"content"`
// 	CreatedAt int64  `json:"created_at"`
// 	UpdatedAt int64  `json:"updated_at"`
// 	DeletedAt int64  `json:"deleted_at"`
// }

// type Message struct {
// 	UID            string `json:"uid"`
// 	Data           string `json:"data"`
// 	ConversationID string `json:"conversation_id"`
// 	RecipientID    string `json:"recipient_id"`
// 	MessageID      string `json:"message_id"`
// 	Category       string `json:"category"`
// 	CreatedAt      int64  `json:"created_at"`
// }

type Snapshot struct {
	SnapshotID string    `bson:"snapshot_id" json:"snapshot_id"`
	RequestID  string    `bson:"request_id" json:"request_id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	AssetID    string    `bson:"asset_id" json:"asset_id"`
	Amount     string    `bson:"amount" json:"amount"`
	Memo       string    `bson:"memo" json:"memo"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

type Refund struct {
	Uid       string    `bson:"uid" json:"uid"`
	Tid       string    `bson:"tid" json:"tid"`
	RequestID string    `bson:"request_id" json:"request_id"`
	Amount    string    `bson:"amount" json:"amount"` // 退款金额
	Action    bool      `bson:"select" json:"select"` // 选择
	Memo      string    `bson:"memo" json:"memo"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// 话题购买历史
type TopicPurchaseHistory struct {
	RequestID  string    `bson:"request_id" json:"request_id"` // 请求ID
	Uid        string    `bson:"uid" json:"uid"`
	Tid        string    `bson:"tid" json:"tid"`
	Action     bool      `bson:"action" json:"action"` // 选择 0 反对, 1 支持
	Amount     string    `bson:"amount" json:"amount"` // 金额
	Memo       string    `bson:"memo" json:"memo"`
	Finished   bool      `bson:"finished" json:"finished"` // 是否完成
	FinishedAt time.Time `bson:"finished_at" json:"finished_at"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

// 用户话题购买
type TopicPurchase struct {
	Uid       string    `bson:"uid" json:"uid"`
	Tid       string    `bson:"tid" json:"tid"`
	YesAmount string    `bson:"yes_amount" json:"yes_amount"` // 支持金额
	NoAmount  string    `bson:"no_amount" json:"no_amount"`   // 反对金额
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (t *Topic) AfterQuery(ctx context.Context) error {
	t.ReadCount++
	return nil
}

// 处理退款记录
type TopicRefundAction struct {
	RequestID string          `json:"request_id"` // 请求ID
	Uid       string          `json:"uid"`
	Tid       string          `json:"tid"`
	Amount    decimal.Decimal `json:"amount"`
	Action    bool            `json:"action"`
	Memo      string          `json:"memo"`
}

type TopicBuyAction struct {
	RequestID string          `json:"request_id"` // 请求ID
	Uid       string          `json:"uid"`
	Tid       string          `json:"tid"`
	Action    bool            `json:"action"`
	Amount    decimal.Decimal `json:"amount"`
	Memo      string          `json:"memo"`
}

type TopicStopAction struct {
	Tid string `json:"tid"`
}

type StopTopicActionResp struct {
	Topic          Topic // 话题
	TopicPurchases []*TopicPurchase
}

type TopicPurchaseRatio struct {
	Uid       string          `json:"uid"`
	WinRatio  string          `json:"win_ratio"`
	WinAmount decimal.Decimal `json:"win_amount"`

	// 购买了的
	YesAmount decimal.Decimal `json:"yes_amount"`
	NoAmount  decimal.Decimal `json:"no_amount"`
}
