package database

import (
	"github.com/lixvyang/betxin.one/internal/model/database/model"
	"github.com/lixvyang/betxin.one/internal/model/database/mysql"
)

type Database interface {
	Init() error
	Close() error
	IUser
}

type IUser interface {
	CheckUser(uid string) int
	GetUserByUid(uid string) (*model.User, int)
	CreateUser(user *model.User) int
	DeleteUser(uid string) int
	UpdateUser(user *model.User) int
	GetUserByFullName(full_name string) (*model.User, int)
}

type ITopic interface {
}

func NewDatabse() Database {
	return mysql.NewMySqlService()
}
