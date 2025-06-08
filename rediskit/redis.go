package rediskit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCtrl struct {
	Rdb *redis.Client
}

func NewRedisCtrl(conf *RedisConf) *RedisCtrl {
	rdb := redis.NewClient(&conf.Options)
	return &RedisCtrl{
		Rdb: rdb,
	}
}

func (s *RedisCtrl) Set(k string, v interface{}, t time.Duration) error {
	cmd := s.Rdb.Set(context.Background(), k, v, t)
	return cmd.Err()
}

func (s *RedisCtrl) Get(k string) (result interface{}, err error) {
	cmd := s.Rdb.Get(context.Background(), k)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	err = cmd.Scan(&result)
	return result, err
}

func (s *RedisCtrl) SAdd(k string, value ...interface{}) (err error) {
	cmd := s.Rdb.SAdd(context.Background(), k, value...)
	return cmd.Err()
}

func (s *RedisCtrl) Exists(key string) (exist bool, err error) {
	exisstint, err := s.Rdb.Exists(context.Background(), key).Result()
	return exisstint == 1, err
}

func (s *RedisCtrl) Expire(key string, d time.Duration) (err error) {
	_, err = s.Rdb.Expire(context.Background(), key, d).Result()
	return err
}
