package cache

import (
	"context"
	"sort"
	"sync"
	"time"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

const defaultMemoryLimitSize = 1000

type memoryValueItem struct {
	lastUsedTime time.Time
	key          string
	value        interface{}
}

type MemoryValue struct {
	// Must not change concurrency with usage
	Name       string // use for log
	MaxSize    int
	CleanCount int

	mu sync.RWMutex
	m  map[string]memoryValueItem
}

func NewMemoryValue(name string) *MemoryValue {
	return &MemoryValue{
		Name:       name,
		m:          make(map[string]memoryValueItem, defaultMemoryLimitSize+1),
		MaxSize:    defaultMemoryLimitSize,
		CleanCount: 300,
	}
}

func (c *MemoryValue) Get(ctx context.Context, key string) (value interface{}, err error) {
	defer func() {
		zc.L(ctx).Debug("Get from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("value", value), zap.Error(err))
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()
	if resp, exist := c.m[key]; exist {
		return resp.value, nil
	}
	return nil, ErrCacheMiss
}

func (c *MemoryValue) Put(ctx context.Context, key string, value interface{}) (err error) {
	defer func() {
		zc.L(ctx).Debug("Put to memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("data_len", value), zap.Error(err))
	}()

	c.mu.Lock()
	c.m[key] = memoryValueItem{key: key, value: value, lastUsedTime: time.Now()}
	if len(c.m) > c.MaxSize {
		go c.clean()
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryValue) Delete(ctx context.Context, key string) (err error) {
	defer func() {
		zc.L(ctx).Debug("Delete from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Error(err))
	}()

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, key)

	return nil
}

func (c *MemoryValue) clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.CleanCount == 0 {
		return
	}

	if len(c.m) <= c.MaxSize {
		return
	}

	if c.CleanCount >= c.MaxSize {
		c.m = make(map[string]memoryValueItem, c.MaxSize+1)
		return
	}

	items := make([]memoryValueItem, 0, len(c.m))
	for _, item := range c.m {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].lastUsedTime.Before(items[j].lastUsedTime)
	})

	for i := 0; i < c.CleanCount; i++ {
		delete(c.m, items[i].key)
	}
}
