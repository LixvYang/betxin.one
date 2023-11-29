package sqlmodel

import (
	"time"

	"github.com/lixvyang/betxin.one/pkg/snowflake"
	"gorm.io/gorm"
)

const TableNameTopic = "topic"

// Topic mapped from table <topic>
type Topic struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"` // ID
	Tid           int64  `gorm:"column:tid;not null" json:"tid"`
	Cid           int64  `gorm:"column:cid;not null;comment:ID" json:"cid"` // ID
	Title         string `gorm:"column:title;not null" json:"title"`
	Intro         string `gorm:"column:intro;not null" json:"intro"`
	Content       string `gorm:"column:content;not null" json:"content"`
	YesRatio      string `gorm:"column:yes_ratio;not null;default:50.00" json:"yes_ratio"`
	NoRatio       string `gorm:"column:no_ratio;not null;default:50.00" json:"no_ratio"`
	YesCount      string `gorm:"column:yes_count;not null;default:0.00000000" json:"yes_count"`
	NoCount       string `gorm:"column:no_count;not null;default:0.00000000" json:"no_count"`
	TotalCount    string `gorm:"column:total_count;not null;default:0.00000000" json:"total_count"`
	CollectCount  int64  `gorm:"column:collect_count;not null" json:"collect_count"`
	ReadCount     int64  `gorm:"column:read_count;not null" json:"read_count"`
	ImgURL        string `gorm:"column:img_url;not null;comment:URL" json:"img_url"` // URL
	IsStop        bool   `gorm:"column:is_stop" json:"is_stop"`
	RefundEndTime int64  `gorm:"column:refund_end_time;not null" json:"refund_end_time"`
	EndTime       int64  `gorm:"column:end_time;not null" json:"end_time"`
	CreatedAt     int64  `gorm:"column:created_at;not null;autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at;not null;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt     int64  `gorm:"column:deleted_at" json:"deleted_at"`
}

func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	t.Tid = snowflake.GenID()
	t.YesRatio = "50.00"
	t.NoRatio = "50.00"
	t.CreatedAt = time.Now().UnixMilli()
	t.UpdatedAt = time.Now().UnixMilli()
	return nil
}

func (t *Topic) AfterFind(tx *gorm.DB) (err error) {
	t.ReadCount++
	return
}
