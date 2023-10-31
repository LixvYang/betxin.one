package redis

type IRedis interface {
	Init() error
	Close() error
}

func NewIRedis() IRedis {
	rds := NewRedisClient()
	return rds
}
