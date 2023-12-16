package core

import (
	"context"
	"time"

	"github.com/pandodao/passport-go/auth"
	"github.com/rs/zerolog"
)

type (
	User struct {
		ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		UID            string    `gorm:"column:uid;not null" json:"uid"`
		IdentityNumber string    `gorm:"column:identity_number;not null" json:"identity_number"`
		FullName       string    `gorm:"column:full_name" json:"full_name"`
		AvatarURL      string    `gorm:"column:avatar_url" json:"avatar_url"`
		SessionID      string    `gorm:"column:session_id" json:"session_id"`
		Biography      string    `gorm:"column:biography" json:"biography"`
		PrivateKey     string    `gorm:"column:private_key" json:"private_key"`
		Contract       string    `gorm:"column:contract" json:"contract"`
		MvmPublicKey   string    `gorm:"column:mvm_public_key" json:"mvm_public_key"`
		IsMvmUser      bool      `gorm:"column:is_mvm_user" json:"is_mvm_user"`
		CreatedAt      time.Time `gorm:"column:created_at;not null" json:"created_at"`
		UpdatedAt      time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	}

	UserStore interface {
		// SELECT
		// 	*
		// FROM user
		// WHERE uid = @uid;
		GetUserByUid(ctx context.Context, uid string) (*User, error)

		// SELECT
		// 	*
		// FROM user
		// WHERE uid in (@uids);
		GetUserByUids(ctx context.Context, uids []string) ([]*User, error)

		// INSERT INTO user
		// 	(
		// 		`full_name`, `avatar_url`,
		// 		`uid`, `identity_number`,
		// 		`session_id`, `is_mvm_user`,
		// 		`biography`, `private_key`,
		// 		`mvm_public_key`, `contract`
		// 	)
		// VALUES
		// 	(
		// 		@user.FullName, @user.AvatarURL,
		// 		@user.UID, @user.IdentityNumber,
		// 		@user.SessionID, @user.IsMvmUser,
		// 		@user.Biography, @user.PrivateKey,
		// 		@user.MvmPublicKey, @user.Contract
		// 	);
		CreateUser(ctx context.Context, user *User) error

		// UPDATE user
		// 	{{set}}
		// 	  `full_name`=@user.FullName,
		// 	  `avatar_url`=@user.AvatarURL,
		// 	  `biography`=@user.Biography
		// 	{{end}}
		// WHERE
		// 	"uid" = @uid;
		UpdateUserInfo(ctx context.Context, uid string, user *User) error
	}

	UserService interface {
		LoginWithMixin(ctx context.Context, logger *zerolog.Logger, authUser *auth.User, isMvmUser bool) (*User, error)
	}
)
