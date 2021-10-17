package lfu

import (
	"github.com/guodf/goutil/cache"
	"sync"
	"sync/atomic"
)

type lfuCache struct {
	cType  cache.CacheType
	lfuMap sync.Map
	size   int64
	min    *lfuCacheItem
}

type lfuCacheItem struct {
	key interface{}
	cache.CacheItem
	times int
	pre   *lfuCacheItem
	next  *lfuCacheItem
}

func (item *lfuCacheItem) fixMinHeap() {
	if item.next == nil {
		return
	}
	if item.times > item.next.times {
		var next *lfuCacheItem
		if item.pre != nil {
			item.pre.next, item.next.pre = item.next, item.pre
			next = item.next
			item.pre = nil
			item.next = nil
		}
		for {
			if item.times > next.times {
				if next.next == nil {
					next.next = item
					item.pre = next
					return
				}
				continue
			}
			next.pre.next, item.pre = item, next.pre
			next.pre, item.next = item, next
		}
	}
}

func newLfuCache() cache.ICache {
	return &lfuCache{
		cType:  cache.LFU,
		lfuMap: sync.Map{},
		size:   0,
	}
}

func (l *lfuCache) Get(key interface{}) interface{} {
	value, _ := l.lfuMap.Load(key)
	item := value.(*lfuCacheItem)
	item.times++
	item.fixMinHeap()
	return item.CacheItem
}

func (l *lfuCache) Set(key interface{}, value interface{}) {
	var lfuItem *lfuCacheItem
	value, ok := l.lfuMap.Load(key)
	cacheItem := value.(cache.CacheItem)
	// 如果key存在则是需要修改value和调整最小堆
	if ok {
		lfuItem = value.(*lfuCacheItem)
		lfuItem.CacheItem = cacheItem
		lfuItem.times++
		min := lfuItem.next
		lfuItem.fixMinHeap()
		if min.key != lfuItem.key {
			l.min = min
		}
	} else { // 新添加的元素需要size+1，将最小堆堆顶指向该元素
		l.checkSize()
		lfuItem = &lfuCacheItem{
			key:       key,
			CacheItem: cacheItem,
			times:     0,
		}
		l.lfuMap.Store(key, lfuItem)
		lfuItem.next = l.min
		l.min = lfuItem
		atomic.AddInt64(&(l.size), 1)
	}
}

func (l *lfuCache) Contains(key interface{}) bool {
	_, ok := l.lfuMap.Load(key)
	return ok
}

func (l *lfuCache) Remove(key interface{}) {
	l.lfuMap.Delete(key)
	atomic.AddInt64(&(l.size), -1)
}

func (l *lfuCache) Size() int64 {
	return l.size
}

func (l *lfuCache) Keys() []interface{} {
	// todo 当获取keys时发生添加、删除的情况如何获取keys
	size := l.size
	keys := make([]interface{}, 0, size)
	var index int64
	l.lfuMap.Range(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return index < size
	})
	return keys
}

func (l *lfuCache) checkSize() {
	if l.size == 1000 {
		l.min = l.min.next
		l.size--
	}
}
