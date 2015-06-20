package cache

import (
	"time"
	"github.com/furio/widserve/cache/local"
)

type CacheGeneric interface {
	Init(config map[string]string)
	Get(key string) (interface{},bool)
	Set(key string, value interface{}, timeout time.Duration)
	Delete(key string)
}

type CacheType int
const (
	Local CacheType = iota
	Redis
)

func GetCacheClient(cacheType CacheType, config map[string]string) CacheGeneric {
	if (cacheType == Local) {
		outCache := local.LocalCache{}
		outCache.Init(config)
		return CacheGeneric(outCache)
	}

	return nil
}