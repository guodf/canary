package lru

import "github.com/guodf/goutil/cache"

type lruCache struct {
	name  string
	cType cache.CacheType
}

type lruCacheItem struct {
}

func newLruCache() cache.ICache {
	return &lruCache{
		name: "LRU",
	}
}

func (lru *lruCache) Get(key interface{}) interface{} {
	return nil
}

func (lru *lruCache) Set(key interface{}, value interface{}) {
	panic("implement me")
}

func (lru *lruCache) Contains(key interface{}) bool {
	panic("implement me")
}

func (lru *lruCache) Remove(key interface{}) {
	panic("implement me")
}

func (lru *lruCache) Size() int64 {
	return 0
}

func (lru *lruCache) Keys() []interface{} {
	panic("implement me")
}
