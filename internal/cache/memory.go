package cache

import (
	"context"
	"github.com/rekby/safemutex"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type MemoryCache struct {
	Name string // use for log

	stateMutex safemutex.RWMutexWithPointers[map[string][]byte]
}

func NewMemoryCache(name string) *MemoryCache {
	return &MemoryCache{
		Name:       name,
		stateMutex: safemutex.RWNewWithPointers(map[string][]byte{}),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (data []byte, err error) {
	defer func() {
		zc.L(ctx).Debug("Get from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Int("data_len", len(data)), zap.Error(err))
	}()

	c.stateMutex.RLock(func(synced map[string][]byte) {
		if resp, exist := synced[key]; exist {
			data = resp
			return
		}
		err = ErrCacheMiss
	})
	return data, err
}

func (c *MemoryCache) Put(ctx context.Context, key string, data []byte) (err error) {
	defer func() {
		zc.L(ctx).Debug("Put to memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Int("data_len", len(data)), zap.Error(err))
	}()

	localCopy := make([]byte, len(data))
	copy(localCopy, data)
	c.stateMutex.Lock(func(synced map[string][]byte) map[string][]byte {
		synced[key] = localCopy
		return synced
	})

	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) (err error) {
	defer func() {
		zc.L(ctx).Debug("Delete from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Error(err))
	}()

	c.stateMutex.Lock(func(synced map[string][]byte) map[string][]byte {
		delete(synced, key)
		return synced
	})

	return nil
}
