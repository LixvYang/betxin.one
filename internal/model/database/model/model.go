package model

import (
	"gorm.io/gorm"
)

type User struct {
	Uid            string `gorm:"type:varchar(36);index;" json:"uid"`
	FullName       string `gorm:"type:varchar(50);not null" json:"full_name"`
	IdentityNumber string `gorm:"varchar(11);not null" json:"identity_number"`
	AvatarUrl      string `gorm:"type:varchar(255);not null" json:"avatar_url"`
	SessionId      string `gorm:"type:varchar(255)" json:"session_id"`
	Biography      string `gorm:"type:varchat(255)" json:"biography"`
	PrivateKey     string `gorm:"type:varchat(255)" json:"private_key"`
	ClientId       string `gorm:"type:varchat(255)" json:"client_id"`
	Contract       string `gorm:"type:varchat(255)" json:"contract"`
	IsMvmUser      int64  `gorm:"type:int(1)" json:"is_mvm_user"`

	gorm.Model
}
