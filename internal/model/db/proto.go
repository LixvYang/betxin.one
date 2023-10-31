package db

import "github.com/lixvyang/betxin.one/internal/model/db/mysql"

type Database interface {
	Init() error
	Close() error
	User
}

type User interface {
	CheckUser(userId string) int
}

func NewDatabse() Database {
	return mysql.NewMySqlService()
}
