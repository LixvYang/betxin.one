package schema

import "time"

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
	YesRatio      string    `bson:"yes_ratio" json:"yes_ratio"`
	NoRatio       string    `bson:"no_ratio" json:"no_ratio"`
	YesCount      string    `bson:"yes_count" json:"yes_count"`
	NoCount       string    `bson:"no_count" json:"no_count"`
	TotalCount    string    `bson:"total_count" json:"total_count"`
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
	UpdatedAt time.Time `json:"updated_at" bson:"created_at"`
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

// type Snapshot struct {
// 	TraceID        string `json:"trace_id"`
// 	Memo           string `json:"memo"`
// 	Type           string `json:"type"`
// 	SnapshotID     string `json:"snapshot_id"`
// 	OpponentID     string `json:"opponent_id"`
// 	AssetID        string `json:"asset_id"`
// 	Amount         string `json:"amount"`
// 	OpeningBalance string `json:"opening_balance"`
// 	ClosingBalance string `json:"closing_balance"`
// 	CreatedAt      string `json:"created_at"`
// }

// type Refund struct {
// 	UID       string `json:"uid"`
// 	AssetID   string `json:"asset_id"`
// 	TraceID   string `json:"trace_id"`
// 	Price     string `json:"price"`  // 退款金额
// 	Select    bool   `json:"select"` // 选择
// 	Memo      string `json:"memo"`
// 	CreatedAt int64  `json:"created_at"`
// }

// type TopicPurchase struct {
// 	UID       string `json:"uid"`
// 	Tid       int64  `json:"tid"`
// 	YesPrice  string `json:"yes_price"` // 支持金额
// 	NoPrice   string `json:"no_price"`  // 反对金额
// 	CreatedAt int64  `json:"created_at"`
// 	UpdatedAt int64  `json:"updated_at"`
// 	DeletedAt int64  `json:"deleted_at"`
// }
