package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	mutex    sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

type Pair struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if cache.capacity == cache.queue.Len() {
		removeEl := cache.queue.Back()
		cache.queue.Remove(removeEl)
		delete(cache.items, removeEl.Value.(*Pair).key)
	}
	_, ok := cache.items[key]
	if ok {
		cache.items[key] = cache.queue.PushFront(&Pair{key, value})
		return true
	}
	queueElem := cache.queue.PushFront(&Pair{key, value})
	cache.items[key] = queueElem
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	val, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(val)
		return val.Value.(*Pair).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue.Clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
