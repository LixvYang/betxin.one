package core

import (
	"time"
)

type (
	Bonuse struct {
		ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		UID       string    `gorm:"column:uid;not null;comment:uid" json:"uid"`          // uid
		Tid       int64     `gorm:"column:tid;not null;comment:id" json:"tid"`           // id
		AssetID   string    `gorm:"column:asset_id;not null;comment:id" json:"asset_id"` // id
		Amount    string    `gorm:"column:amount;not null" json:"amount"`
		Memo      string    `gorm:"column:memo;not null" json:"memo"`
		TraceID   string    `gorm:"column:trace_id;not null" json:"trace_id"`
		CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
		DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	}

	BonuseStore interface {
	}
)
