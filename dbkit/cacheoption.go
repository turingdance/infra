package dbkit

import (
	"github.com/turingdance/gormcache/serializer"
	"github.com/turingdance/gormcache/store"
)

const (
	LevelDisable = 0 //禁止
	LevelModel   = 1 //只缓存模型
	LevelSearch  = 2 //查询缓存
	MaxExpires   = 43200
	MinExpires   = 30
)

type CacheOption struct {
	store      store.Store
	serializer serializer.Serializer
}
type CacheOptionSetting func(*CacheOption)

func UseStore(store store.Store) CacheOptionSetting {
	return func(co *CacheOption) {
		co.store = store
	}
}
func UseSerializer(serializer serializer.Serializer) CacheOptionSetting {
	return func(co *CacheOption) {
		co.serializer = serializer
	}
}
func NewCacheOption(settings ...CacheOptionSetting) *CacheOption {
	ret := &CacheOption{}
	for _, v := range settings {
		v(ret)
	}
	return ret
}
func UseCache(opt *CacheOption) Option {
	return func(dc *DbContext) {
		dc.cacheOption = opt
	}
}
