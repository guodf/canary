package fifo

import (
	"github.com/guodf/goutil/cache"
	"sync"
)

type fifoCache struct {
	first   *fifoCacheItem
	tail    *fifoCacheItem
	fifoMap sync.Map
	size    int64
}

type fifoCacheItem struct {
	key interface{}
	cache.CacheItem
	pre  *fifoCacheItem
	next *fifoCacheItem
}

func newFifoCache() cache.ICache {
	return &fifoCache{
		fifoMap: sync.Map{},
	}
}

func (f *fifoCache) Get(key interface{}) interface{} {
	_, value := f.fifoMap.Load(key)
	return value
}

func (f *fifoCache) Set(key interface{}, value interface{}) {
	var cacheItem fifoCacheItem
	if f.size == 0 {
		cacheItem = fifoCacheItem{
			key:       key,
			CacheItem: value.(cache.CacheItem),
		}
		f.first = &cacheItem
		f.tail = &cacheItem
		f.fifoMap.Store(key, &cacheItem)
		f.size++
		return
	}
	v, ok := f.fifoMap.Load(key)
	if ok {
		v.(*fifoCacheItem).CacheItem = value.(cache.CacheItem)
		return
	}
	f.checkSize()
	cacheItem = fifoCacheItem{
		key:       key,
		CacheItem: value.(cache.CacheItem),
	}
	f.fifoMap.Store(key, value)

}

func (f *fifoCache) Contains(key interface{}) bool {
	_, ok := f.fifoMap.Load(key)
	return ok
}

func (f *fifoCache) Remove(key interface{}) {
	v, ok := f.fifoMap.Load(key)
	if !ok {
		return
	}
	cacheItem := v.(*fifoCacheItem)
	if cacheItem.pre == nil {
		f.removeFirst()
		return
	}
	if cacheItem.next == nil {
		f.removeTail()
		return
	}
	cacheItem.pre.next, cacheItem.next.pre = cacheItem.next, cacheItem.pre
	f.fifoMap.Delete(key)

}

func (f *fifoCache) Size() int64 {
	return f.size
}

func (f *fifoCache) Keys() []interface{} {
	// todo 当获取keys时发生添加、删除的情况如何获取keys
	size := f.size
	keys := make([]interface{}, 0, size)
	var index int64
	f.fifoMap.Range(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return index < size
	})
	return keys
}

func (f *fifoCache) checkSize() {
	if f.size == 1000 {
		f.size--
		f.first = f.first.next
		f.first.pre = nil
		f.fifoMap.Delete(f.first.key)
	}
}

func (f *fifoCache) removeFirst() {
	f.fifoMap.Delete(f.first.key)
	f.first = f.first.next
	f.size--
	if f.first == nil {
		f.tail = nil
	} else {
		f.first.pre = nil
	}
}

func (f *fifoCache) removeTail() {
	f.fifoMap.Delete(f.tail.key)
	f.tail = f.tail.pre
	f.size--
	if f.tail == nil {
		f.first = nil
	} else {
		f.tail.next = nil
	}
}
