package substack

import (
	"sync"
	"time"
)

type cacheEntry struct {
	posts []Post
	exp   time.Time
}

type FeedCache struct {
	mu     sync.RWMutex
	store  map[string]cacheEntry // key: url|limit
	ttl    time.Duration
}

func NewFeedCache(ttl time.Duration) *FeedCache {
	return &FeedCache{store: make(map[string]cacheEntry), ttl: ttl}
}

func (c *FeedCache) key(url string, limit int) string {
	return url + "|" + time.Duration(limit).String()
}

func (c *FeedCache) Get(url string, limit int) ([]Post, bool) {
	k := c.key(url, limit)
	c.mu.RLock()
	defer c.mu.RUnlock()
	ce, ok := c.store[k]
	if !ok || time.Now().After(ce.exp) {
		return nil, false
	}
	return ce.posts, true
}

func (c *FeedCache) Set(url string, limit int, posts []Post) {
	k := c.key(url, limit)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[k] = cacheEntry{posts: posts, exp: time.Now().Add(c.ttl)}
}

func (c *FeedCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]cacheEntry)
}
