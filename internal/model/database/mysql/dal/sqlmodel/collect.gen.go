// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package sqlmodel

const TableNameCollect = "collect"

// Collect mapped from table <collect>
type Collect struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UID       string `gorm:"column:uid;not null" json:"uid"`
	Tid       int64  `gorm:"column:tid;not null" json:"tid"`
	Status    bool   `gorm:"column:status;not null;comment:状态" json:"status"` // 状态
	CreatedAt int64  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName Collect's table name
func (*Collect) TableName() string {
	return TableNameCollect
}
