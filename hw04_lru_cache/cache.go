package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := cache.items[key]
	isNewItem := item == nil

	newCacheItem := cacheItem{
		value: value,
		key:   key,
	}

	if isNewItem {
		cache.queue.PushFront(newCacheItem)
		cache.items[key] = cache.queue.Front()

		if cache.queue.Len() > cache.capacity {
			lastItem := cache.queue.Back()

			cache.queue.Remove(lastItem)
			cache.items[lastItem.Value.(cacheItem).key] = nil
		}
	} else {
		cache.items[key].Value = newCacheItem
		cache.queue.MoveToFront(cache.items[key])
	}

	return !isNewItem
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := cache.items[key]

	if isNewItem := item == nil; isNewItem {
		return nil, false
	}

	cache.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, true
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}
