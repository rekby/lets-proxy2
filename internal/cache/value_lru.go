package cache

import (
	"context"
	"github.com/rekby/safemutex"
	"math"
	"sort"
	"sync/atomic"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

const defaultMemoryLimitSize = 1000
const defaultLRUCleanCount = 300

type memoryValueLRUItem struct {
	key   string
	value interface{}

	lastUsedTime atomic.Uint64
}

type MemoryValueLRU struct {
	// Must not change concurrency with usage
	Name       string // use for log
	MaxSize    int
	CleanCount int

	lastTime uint64
	mu       safemutex.RWMutexWithPointers[memoryValueLRUSynced]
}

type memoryValueLRUSynced struct {
	Items map[string]*memoryValueLRUItem // always stored non nil items
}

func NewMemoryValueLRU(name string) *MemoryValueLRU {
	return &MemoryValueLRU{
		Name: name,
		mu: safemutex.RWNewWithPointers(memoryValueLRUSynced{
			Items: make(map[string]*memoryValueLRUItem, defaultMemoryLimitSize+1),
		}),
		MaxSize:    defaultMemoryLimitSize,
		CleanCount: defaultLRUCleanCount,
	}
}

func (c *MemoryValueLRU) Get(ctx context.Context, key string) (value interface{}, err error) {
	defer func() {
		zc.L(ctx).Debug("Get from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("value", value), zap.Error(err))
	}()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		if resp, exist := synced.Items[key]; exist {
			resp.lastUsedTime.Store(c.time())
			value = resp.value
			return
		}
		err = ErrCacheMiss
	})

	return value, err
}

func (c *MemoryValueLRU) Put(ctx context.Context, key string, value interface{}) (err error) {
	defer func() {
		zc.L(ctx).Debug("Put to memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("data_len", value), zap.Error(err))
	}()

	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items[key] = &memoryValueLRUItem{key: key, value: value, lastUsedTime: newUint64Atomic(c.time())}
		if len(synced.Items) > c.MaxSize {
			// handlepanic: no external call
			go c.clean()
		}

		return synced
	})

	return nil
}

func (c *MemoryValueLRU) Delete(ctx context.Context, key string) (err error) {
	defer func() {
		zc.L(ctx).Debug("Delete from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Error(err))
	}()

	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		delete(synced.Items, key)
		return synced
	})

	return nil
}

func (c *MemoryValueLRU) time() uint64 {
	res := atomic.AddUint64(&c.lastTime, 1)
	if res == math.MaxUint64/2 {
		// handlepanic: no external call
		go c.renumberTime()
	}
	return res
}

func (c *MemoryValueLRU) renumberTime() {
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		items := c.getSortedItems(&synced)
		for i, item := range items {
			item.lastUsedTime.Store(uint64(i))
		}
		return synced
	})
}

// must called from locked state
func (c *MemoryValueLRU) getSortedItems(synced *memoryValueLRUSynced) []*memoryValueLRUItem {
	items := make([]*memoryValueLRUItem, 0, len(synced.Items))
	for _, item := range synced.Items {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].lastUsedTime.Load() < items[j].lastUsedTime.Load()
	})
	return items
}

func (c *MemoryValueLRU) clean() {
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		if c.CleanCount == 0 {
			return synced
		}

		if len(synced.Items) <= c.MaxSize {
			return synced
		}

		if c.CleanCount >= c.MaxSize {
			synced.Items = make(map[string]*memoryValueLRUItem, c.MaxSize+1)
			return synced
		}

		items := c.getSortedItems(&synced)

		for i := 0; i < c.CleanCount; i++ {
			delete(synced.Items, items[i].key)
		}

		return synced
	})
}

func newUint64Atomic(val uint64) atomic.Uint64 {
	var v atomic.Uint64
	v.Store(val)
	return v
}
