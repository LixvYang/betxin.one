package schema

type User struct {
	IdentityNumber string `json:"identity_number"`
	UID            string `json:"uid"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
	SessionID      string `json:"session_id"`
	Biography      string `json:"biography"`
	PrivateKey     string `json:"private_key"`
	ClientID       string `json:"client_id"`
	Contract       string `json:"contract"`
	IsMvmUser      bool   `json:"is_mvm_user"`
}

type Topic struct {
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

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Bonuse struct {
	UID       string `json:"uid"`      // uid
	Tid       int64  `json:"tid"`      // id
	AssetID   string `json:"asset_id"` // id
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
	TraceID   string `json:"trace_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

type Collect struct {
	UID       string `json:"uid"`
	Tid       int64  `json:"tid"`
	Status    bool   `json:"status"` // 状态
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Feedback struct {
	UID       string `json:"uid"`
	Fid       string `json:"fid"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

type Message struct {
	UID            string `json:"uid"`
	Data           string `json:"data"`
	ConversationID string `json:"conversation_id"`
	RecipientID    string `json:"recipient_id"`
	MessageID      string `json:"message_id"`
	Category       string `json:"category"`
	CreatedAt      int64  `json:"created_at"`
}

type Snapshot struct {
	TraceID        string `json:"trace_id"`
	Memo           string `json:"memo"`
	Type           string `json:"type"`
	SnapshotID     string `json:"snapshot_id"`
	OpponentID     string `json:"opponent_id"`
	AssetID        string `json:"asset_id"`
	Amount         string `json:"amount"`
	OpeningBalance string `json:"opening_balance"`
	ClosingBalance string `json:"closing_balance"`
	CreatedAt      string `json:"created_at"`
}

type Refund struct {
	UID       string `json:"uid"`
	AssetID   string `json:"asset_id"`
	TraceID   string `json:"trace_id"`
	Price     string `json:"price"`  // 退款金额
	Select    bool   `json:"select"` // 选择
	Memo      string `json:"memo"`
	CreatedAt int64  `json:"created_at"`
}

type TopicPurchase struct {
	UID       string `json:"uid"`
	Tid       int64  `json:"tid"`
	YesPrice  string `json:"yes_price"` // 支持金额
	NoPrice   string `json:"no_price"`  // 反对金额
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}
