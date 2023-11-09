package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/configs"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var (
	Nil = redis.Nil
)

type Cache struct {
	cli *redis.Client
}

func New(conf *configs.RedisConfig) *Cache {
	cache := &Cache{}
	cache.cli = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	})

	res, err := cache.cli.Ping(context.Background()).Result()
	if err != nil {
		logger.Lg.Panic().Err(err).Msg("redis ping error")
		panic(err)
	}
	logger.Lg.Info().Str("redis", res).Send()
	return cache
}

func (r *Cache) Close() error {
	return r.cli.Close()
}

func (r *Cache) Get(ctx context.Context, key string) (ret []byte, err error) {
	ret, err = r.cli.Get(ctx, key).Bytes()
	return ret, err
}

func (r *Cache) HGet(ctx context.Context, key string, field string) (ret []byte, err error) {
	return r.cli.HGet(ctx, key, field).Bytes()
}

func (r *Cache) GetRes(ctx context.Context, key string, res any) (err error) {
	ret, err := r.cli.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(ret, &res); err != nil {
		return err
	}
	return nil
}

func (r *Cache) Set(ctx context.Context, key string, value []byte, expireSeconds int32) (err error) {
	return r.cli.Set(ctx, key, value, time.Duration(expireSeconds)*time.Second).Err()
}

func (r *Cache) HSet(ctx context.Context, key string, field string, value []byte) (err error) {
	return r.cli.HSet(ctx, key, field, value).Err()
}

func (r *Cache) HDel(ctx context.Context, key string, field string) (err error) {
	return r.cli.HDel(ctx, key, field).Err()
}

func (m *Cache) Delete(ctx context.Context, key string) (err error) {
	return m.cli.Del(ctx, key).Err()
}
