package db

type Database interface {
	Init() error
	Close() error
}
