package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/internal/utils/convert"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	Nil = redis.Nil
)

type Cache struct {
	cli *redis.Client
}

func New(logger *zerolog.Logger, conf *config.RedisConfig) *Cache {
	return nil
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
		logger.Panic().Err(err).Msg("redis ping error")
		panic(err)
	}
	logger.Info().Str("redis", res).Send()
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

	if err = convert.Unmarshal(ret, &res); err != nil {
		return err
	}
	return nil
}

func (r *Cache) HGetAll(ctx context.Context, key string, field string) (map[string]string, error) {
	return r.cli.HGetAll(ctx, key).Result()
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

func (m *Cache) ZAdd(ctx context.Context, key string, score float64, val string) (err error) {
	return m.cli.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: val,
	}).Err()
}

func (m *Cache) ZRangeByScore(ctx context.Context, key string, min, max string, offset, limit int64) (ret []string, err error) {
	return m.cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  limit,
	}).Result()
}
