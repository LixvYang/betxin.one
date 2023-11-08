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
	Tid           string `json:"tid"`
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
	RefundEndTime int64  `json:"refund_end_time"`
	EndTime       int64  `json:"end_time"`
}
