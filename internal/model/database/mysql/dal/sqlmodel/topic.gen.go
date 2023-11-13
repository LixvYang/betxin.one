package sqlmodel

import (
	"time"

	"github.com/lixvyang/betxin.one/pkg/snowflake"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const TableNameTopic = "topic"

// Topic mapped from table <topic>
type Topic struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:话题自增ID" json:"id"` // 话题自增ID
	Tid           int64  `gorm:"column:tid;not null" json:"tid"`
	Cid           int64  `gorm:"column:cid;not null;comment:分类ID" json:"cid"`                                   // 分类ID
	Title         string `gorm:"column:title;not null;comment:标题" json:"title"`                                 // 标题
	Intro         string `gorm:"column:intro;not null;comment:概述" json:"intro"`                                 // 概述
	Content       string `gorm:"column:content;not null;comment:内容" json:"content"`                             // 内容
	YesRatio      string `gorm:"column:yes_ratio;not null;default:50.00;comment:赞成率" json:"yes_ratio"`          // 赞成率
	NoRatio       string `gorm:"column:no_ratio;not null;default:50.00;comment:反对率" json:"no_ratio"`            // 反对率
	YesCount      string `gorm:"column:yes_count;not null;default:0.00000000;comment:赞成计数" json:"yes_count"`    // 赞成计数
	NoCount       string `gorm:"column:no_count;not null;default:0.00000000;comment:反对计数" json:"no_count"`      // 反对计数
	TotalCount    string `gorm:"column:total_count;not null;default:0.00000000;comment:总计数" json:"total_count"` // 总计数
	CollectCount  int64  `gorm:"column:collect_count;not null;comment:收藏数" json:"collect_count"`                // 收藏数
	ReadCount     int64  `gorm:"column:read_count;not null;comment:阅读数" json:"read_count"`                      // 阅读数
	ImgURL        string `gorm:"column:img_url;not null;comment:图片URL" json:"img_url"`                          // 图片URL
	IsStop        bool   `gorm:"column:is_stop;comment:是否结束" json:"is_stop"`                                    // 是否结束
	RefundEndTime int64  `gorm:"column:refund_end_time;not null;comment:退款截止时间" json:"refund_end_time"`         // 退款截止时间
	EndTime       int64  `gorm:"column:end_time;not null;comment:话题结束时间" json:"end_time"`                       // 话题结束时间
	IsDeleted     bool   `gorm:"column:is_deleted;not null;comment:是否删除" json:"is_deleted"`                     // 是否删除
	CreatedAt     int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`                     // 创建时间
	UpdatedAt     int64  `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"`                     // 更新时间
	DeletedAt     int64  `gorm:"column:deleted_at;comment:删除时间" json:"deleted_at"`                              // 删除时间
}

// TableName Topic's table name
func (*Topic) TableName() string {
	return TableNameTopic
}

func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	t.Tid = snowflake.GenID()
	t.YesRatio = "50.00"
	t.NoRatio = "50.00"
	return nil
}

func (t *Topic) BeforeUpdate(tx *gorm.DB) error {
	if t.IsStop || time.Now().After(time.UnixMicro(t.EndTime)) {
		return errors.New("topic already stop")
	}
	decimal.DivisionPrecision = 2
	yesCnt, _ := decimal.NewFromString(t.YesCount)
	totalCnt, err := decimal.NewFromString(t.TotalCount)
	if err != nil {
		return err
	}
	yesRatio := yesCnt.Div(totalCnt)
	t.YesRatio = yesRatio.String()
	t.NoRatio = decimal.NewFromInt(100).Sub(yesRatio).String()
	return nil
}

func (t *Topic) AfterFind(tx *gorm.DB) (err error) {
	t.ReadCount++
	return
}
