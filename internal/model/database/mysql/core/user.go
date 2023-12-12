package core

import (
	"context"
)

type (
	User struct {
		ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		IdentityNumber string `gorm:"column:identity_number;not null" json:"identity_number"`
		UID            string `gorm:"column:uid;not null" json:"uid"`
		FullName       string `gorm:"column:full_name" json:"full_name"`
		AvatarURL      string `gorm:"column:avatar_url" json:"avatar_url"`
		SessionID      string `gorm:"column:session_id" json:"session_id"`
		Biography      string `gorm:"column:biography" json:"biography"`
		PrivateKey     string `gorm:"column:private_key" json:"private_key"`
		ClientID       string `gorm:"column:client_id" json:"client_id"`
		Contract       string `gorm:"column:contract" json:"contract"`
		IsMvmUser      bool   `gorm:"column:is_mvm_user" json:"is_mvm_user"`
		CreatedAt      int64  `gorm:"column:created_at;not null" json:"created_at"`
		UpdatedAt      int64  `gorm:"column:updated_at;not null" json:"updated_at"`
	}

	UserStore interface {
		// SELECT
		// 	*
		// FROM @@table
		// WHERE uid = @uid;
		GetUserByUid(ctx context.Context, uid string) (*User, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE uid in (@uids);
		GetUserByUids(ctx context.Context, uids []string) ([]*User, error)

		// INSERT INTO user
		// 	(
		// 		"full_name", "avatar_url",
		// 		"uid", "identity_number",
		// 		"session_id", "is_mvm_user",
		// 		"biography", "private_key",
		// 		"client_id", "contract"
		// 	)
		// VALUES
		// 	(
		// 		@user.FullName, @user.AvatarURL,
		// 		@user.UID, @user.IdentityNumber,
		// 		@user.Biography, @user.PrivateKey,
		// 		@user.ClientID, @user.Contract
		// 	)
		// RETURNING uid;
		CreateUser(ctx context.Context, user *User) (string, error)

		// UPDATE @@table
		// 	{{set}}
		// 	  "full_name"=@user.FullName,
		// 	  "avatar_url"=@user.AvatarURL,
		// 	  "biography"=@user.Biography
		// 	{{end}}
		// WHERE
		// 	"uid" = @uid;
		UpdateUserInfo(ctx context.Context, uid string, user *User) error
	}
)
