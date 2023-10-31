package redis

import (
	"context"
	"fmt"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var (
	Nil = redis.Nil
)

type RedisClient struct {
	rds *redis.Client
}

func NewRedisClient() *RedisClient {
	rds := &RedisClient{}
	if err := rds.Init(); err != nil {
		logger.Lg.Error().Err(err).Msg("[NewRedisClient][panic]")
		panic(err)
	}
	return rds
}

func (r *RedisClient) Init() error {
	r.rds = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", configs.Conf.RedisConfig.Host, configs.Conf.RedisConfig.Port),
		Password:     configs.Conf.RedisConfig.Password, // no password set
		DB:           configs.Conf.RedisConfig.DB,       // use default DB
		PoolSize:     configs.Conf.RedisConfig.PoolSize,
		MinIdleConns: configs.Conf.RedisConfig.MinIdleConns,
	})

	res, err := r.rds.Ping(context.Background()).Result()
	if err != nil {
		logger.Lg.Panic().Err(err).Msg("redis ping error")
		return err
	}
	logger.Lg.Info().Str("redis", res).Send()
	return nil
}

func (r *RedisClient) Close() error {
	return r.rds.Close()
}
