package cache

type CacheType int

const (
	LFU CacheType = iota + 1
	LRU
	LRU_K
	TwoQueue
	ARC
	FIFO
)

type ICache interface {
	Get(key interface{}) interface{}
	Set(key interface{}, value interface{})
	Contains(key interface{}) bool
	Remove(key interface{})
	Size() int64
	Keys() []interface{}
}

type CachePolicy struct {
	PolicyType int
	Content    interface{}
}

type cache struct {
	name   string
	policy CachePolicy
}

type cacheKey struct {
	key interface{}
}

type CacheItem struct {
	key   cacheKey
	value interface{}
}

func NewCache(name string, policy CachePolicy) cache {
	return cache{name: name, policy: policy}
}
