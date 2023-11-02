package mysql

import (
	"github.com/lixvyang/betxin.one/internal/model/database/model"
	"github.com/lixvyang/betxin.one/internal/utils/errmsg"
	"gorm.io/gorm"
)

type UserModel struct {
	UserDB *gorm.DB
}

func NewUserModel(db *gorm.DB) UserModel {
	return UserModel{
		db.Model(&model.User{}),
	}
}

// CheckUser 查询用户是否存在
func (um *UserModel) CheckUser(userId string) int {
	var user model.User
	um.UserDB.Where("mixin_uuid = ?", userId).Last(&user)
	if user.ID == 0 {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// Create user
func (um *UserModel) CreateUser(user *model.User) int {
	if err := um.UserDB.Create(&user).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func (um *UserModel) GetUserByUid(usedId string) (*model.User, int) {
	var user *model.User
	if err := um.UserDB.Model(&user).Where("mixin_uuid = ?", usedId).First(&user).Error; err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// Delete user
func (um *UserModel) DeleteUser(userId string) int {
	if err := um.UserDB.Where("user_id = ?", userId).Delete(model.User{}).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// EditUser 编辑用户信息
func (um *UserModel) UpdateUser(user *model.User) int {
	tx := um.UserDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errmsg.ERROR
	}

	// 锁住指定 id 的 User 记录
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", user.Uid).Error; err != nil {
		tx.Rollback()
		return errmsg.ERROR
	}

	var maps = make(map[string]interface{})
	maps["full_name"] = user.FullName
	if err := um.UserDB.Model(&model.User{}).Where("identity_number = ? ", user.IdentityNumber).Updates(maps).Error; err != nil {
		return errmsg.ERROR
	}
	if err := tx.Commit().Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetUserByName gets an user by the username.
func (um *UserModel) GetUserByFullName(full_name string) (*model.User, int) {
	var user *model.User
	if err := um.UserDB.Where("full_name LIKE ?", full_name+"%").First(&user).Error; err != nil {
		return nil, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// List users
func (um *UserModel) ListUser(offset, limit int) ([]model.User, int, int) {
	var users []model.User
	var count int64
	if err := um.UserDB.Model(&model.User{}).Count(&count).Error; err != nil {
		return users, 0, errmsg.ERROR
	}

	if err := um.UserDB.Where("").Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, 0, errmsg.ERROR
	}
	return users, int(count), errmsg.SUCCSE
}
