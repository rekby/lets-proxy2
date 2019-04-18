package cache

import (
	"context"
	"sync"

	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type MemoryCache struct {
	Name string // use for log

	mu sync.RWMutex
	m  map[string][]byte
}

func NewMemoryCache(name string) *MemoryCache {
	return &MemoryCache{
		Name: name,
		m:    make(map[string][]byte),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (data []byte, err error) {
	defer func() {
		zc.L(ctx).Debug("Get from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Int("data_len", len(data)), zap.Error(err))
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()
	if resp, exist := c.m[key]; exist {
		return resp, nil
	} else {
		return nil, ErrCacheMiss
	}
}

func (c *MemoryCache) Put(ctx context.Context, key string, data []byte) (err error) {
	defer func() {
		zc.L(ctx).Debug("Put to memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Int("data_len", len(data)), zap.Error(err))
	}()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = data
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) (err error) {
	defer func() {
		zc.L(ctx).Debug("Delete from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Error(err))
	}()

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, key)

	return nil
}
