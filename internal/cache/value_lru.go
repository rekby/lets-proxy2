package cache

import (
	"context"
	"math"
	"sort"
	"sync"
	"sync/atomic"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

const defaultMemoryLimitSize = 1000
const defaultLRUCleanCount = 300

type memoryValueLRUItem struct {
	key   string
	value interface{}

	m            sync.Mutex // sync update lastUsedTime in Get method
	lastUsedTime uint64
}

type MemoryValueLRU struct {
	// Must not change concurrency with usage
	Name       string // use for log
	MaxSize    int
	CleanCount int

	lastTime uint64
	mu       sync.RWMutex
	m        map[string]*memoryValueLRUItem // stored always non nil item
}

func NewMemoryValueLRU(name string) *MemoryValueLRU {
	return &MemoryValueLRU{
		Name:       name,
		m:          make(map[string]*memoryValueLRUItem, defaultMemoryLimitSize+1),
		MaxSize:    defaultMemoryLimitSize,
		CleanCount: defaultLRUCleanCount,
	}
}

func (c *MemoryValueLRU) Get(ctx context.Context, key string) (value interface{}, err error) {
	defer func() {
		zc.L(ctx).Debug("Get from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("value", value), zap.Error(err))
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()

	if resp, exist := c.m[key]; exist {
		resp.m.Lock()
		resp.lastUsedTime = c.time()
		resp.m.Unlock()
		return resp.value, nil
	}
	return nil, ErrCacheMiss
}

func (c *MemoryValueLRU) Put(ctx context.Context, key string, value interface{}) (err error) {
	defer func() {
		zc.L(ctx).Debug("Put to memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Reflect("data_len", value), zap.Error(err))
	}()

	c.mu.Lock()
	c.m[key] = &memoryValueLRUItem{key: key, value: value, lastUsedTime: c.time()}
	if len(c.m) > c.MaxSize {
		// handlepanic: no external call
		go c.clean()
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryValueLRU) Delete(ctx context.Context, key string) (err error) {
	defer func() {
		zc.L(ctx).Debug("Delete from memory cache", zap.String("cache_name", c.Name),
			zap.String("key", key), zap.Error(err))
	}()

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, key)

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
	c.mu.Lock()

	items := c.getSortedItems()
	for i, item := range items {
		item.lastUsedTime = uint64(i)
	}

	c.mu.Unlock()
}

// must called from locked state
func (c *MemoryValueLRU) getSortedItems() []*memoryValueLRUItem {
	items := make([]*memoryValueLRUItem, 0, len(c.m))
	for _, item := range c.m {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].lastUsedTime < items[j].lastUsedTime
	})
	return items
}

func (c *MemoryValueLRU) clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.CleanCount == 0 {
		return
	}

	if len(c.m) <= c.MaxSize {
		return
	}

	if c.CleanCount >= c.MaxSize {
		c.m = make(map[string]*memoryValueLRUItem, c.MaxSize+1)
		return
	}

	items := c.getSortedItems()

	for i := 0; i < c.CleanCount; i++ {
		delete(c.m, items[i].key)
	}
}
