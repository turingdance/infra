package cachekit

import (
	"context"
	"time"

	"github.com/turingdance/infra/rediskit"
)

type rediscache struct {
	ctrl *rediskit.RedisCtrl
}

func NewRedisCache(conf *rediskit.RedisConf) *rediscache {
	ctrl := rediskit.NewRedisCtrl(conf)
	return &rediscache{
		ctrl: ctrl,
	}
}
func (r *rediscache) Set(k string, v interface{}, d time.Duration) error {
	cmd := r.ctrl.Rdb.Set(context.Background(), k, v, d)
	return cmd.Err()
}

// 获得
func (r *rediscache) Get(k string) (result interface{}, err error) {
	return r.ctrl.Rdb.Get(context.Background(), k).Result()
}
