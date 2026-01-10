package cache

import (
	"context"
	"sync"
	"time"
)

// CacheItem представляет элемент кеша с TTL
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// SimpleCache - простой потокобезопасный кеш в памяти с TTL
type SimpleCache struct {
	mu     sync.RWMutex
	items  map[string]CacheItem
	ttl    time.Duration
	ctx    context.Context
	cancel context.CancelFunc
}

// NewSimpleCache создает новый кеш с указанным TTL
func NewSimpleCache(ttl time.Duration) *SimpleCache {
	ctx, cancel := context.WithCancel(context.Background())
	cache := &SimpleCache{
		items:  make(map[string]CacheItem),
		ttl:    ttl,
		ctx:    ctx,
		cancel: cancel,
	}

	// Запускаем фоновую очистку устаревших элементов каждые 5 минут
	go cache.cleanupExpired()

	return cache
}

// Set добавляет элемент в кеш
func (c *SimpleCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

// Get получает элемент из кеша
// Возвращает (value, true) если найден и не истек
// Возвращает (nil, false) если не найден или истек
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Проверяем истек ли элемент
	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Delete удаляет элемент из кеша
func (c *SimpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear очищает весь кеш
func (c *SimpleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
}

// Size возвращает количество элементов в кеше
func (c *SimpleCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// cleanupExpired удаляет истекшие элементы из кеша
func (c *SimpleCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, item := range c.items {
				if now.After(item.Expiration) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// GetStats возвращает статистику кеша
func (c *SimpleCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	valid := 0
	expired := 0
	now := time.Now()

	for _, item := range c.items {
		if now.After(item.Expiration) {
			expired++
		} else {
			valid++
		}
	}

	return map[string]interface{}{
		"total_items":   len(c.items),
		"valid_items":   valid,
		"expired_items": expired,
		"ttl_seconds":   c.ttl.Seconds(),
	}
}

// Close останавливает фоновую очистку кеша
func (c *SimpleCache) Close() {
	c.cancel()
}
