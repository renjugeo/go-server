package localcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type LocalCacheProvider interface {
	Set(key string, value interface{}, d time.Duration)
	Get(key string) (interface{}, bool)
}

func NewLocalCacheProvider(expiration time.Duration) LocalCacheProvider {
	return cache.New(expiration, 2*expiration)
}
