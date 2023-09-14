package cache

import (
	"math"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestValueLRUAsCache(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	c := NewMemoryValueLRU("test")
	res, err := c.Get(ctx, "asd")
	e.Nil(res)
	e.CmpDeeply(err, ErrCacheMiss)

	data := []byte("aaa")
	err = c.Put(ctx, "asd", data)
	e.CmpNoError(err)

	res, err = c.Get(ctx, "asd")
	e.CmpDeeply(res, data)
	e.CmpNoError(err)

	err = c.Delete(ctx, "asd")
	e.CmpNoError(err)

	err = c.Delete(ctx, "non-existed-key")
	e.CmpNoError(err)

	res, err = c.Get(ctx, "asd")
	e.Nil(res)
	e.CmpDeeply(err, ErrCacheMiss)
}

func TestValueLRULimitAtPut(t *testing.T) {
	td := testdeep.NewT(t)

	ctx, flush := th.TestContext(t)
	defer flush()

	wait := func() { time.Sleep(time.Millisecond * 10) }

	var c *MemoryValueLRU
	var res interface{}
	var err error

	c = NewMemoryValueLRU("test")
	c.MaxSize = 5
	c.CleanCount = 3

	err = c.Put(ctx, "1", 1)
	td.CmpNoError(err)
	err = c.Put(ctx, "2", 2)
	td.CmpNoError(err)
	err = c.Put(ctx, "3", 3)
	td.CmpNoError(err)
	err = c.Put(ctx, "4", 4)
	td.CmpNoError(err)
	err = c.Put(ctx, "5", 5)
	td.CmpNoError(err)

	res, err = c.Get(ctx, "1")
	td.CmpDeeply(res, 1)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "2")
	td.CmpDeeply(res, 2)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "3")
	td.CmpDeeply(res, 3)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "4")
	td.CmpDeeply(res, 4)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "5")
	td.CmpDeeply(res, 5)
	td.CmpNoError(err)

	err = c.Put(ctx, "6", 6)
	td.CmpNoError(err)
	wait()

	res, err = c.Get(ctx, "1")
	td.Nil(res)
	td.CmpDeeply(err, ErrCacheMiss)
	res, err = c.Get(ctx, "2")
	td.Nil(res)
	td.CmpDeeply(err, ErrCacheMiss)
	res, err = c.Get(ctx, "3")
	td.Nil(res)
	td.CmpDeeply(err, ErrCacheMiss)
	res, err = c.Get(ctx, "4")
	td.CmpDeeply(res, 4)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "5")
	td.CmpDeeply(res, 5)
	td.CmpNoError(err)
	res, err = c.Get(ctx, "6")
	td.CmpDeeply(res, 6)
	td.CmpNoError(err)
}

func TestValueLRULimitClean(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var c = NewMemoryValueLRU("test")

	c.MaxSize = 5
	c.CleanCount = 0
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(6)}
		return synced
	})

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 6)
		td.CmpDeeply(synced.Items["1"].value, 1)
		td.CmpDeeply(synced.Items["2"].value, 2)
		td.CmpDeeply(synced.Items["3"].value, 3)
		td.CmpDeeply(synced.Items["4"].value, 4)
		td.CmpDeeply(synced.Items["5"].value, 5)
		td.CmpDeeply(synced.Items["6"].value, 6)
	})

	c.MaxSize = 5
	c.CleanCount = 3
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		return synced
	})

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 5)
		td.CmpDeeply(synced.Items["1"].value, 1)
		td.CmpDeeply(synced.Items["2"].value, 2)
		td.CmpDeeply(synced.Items["3"].value, 3)
		td.CmpDeeply(synced.Items["4"].value, 4)
		td.CmpDeeply(synced.Items["5"].value, 5)
	})

	c.MaxSize = 5
	c.CleanCount = 2

	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(6)}
		return synced
	})

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 4)
		td.Nil(synced.Items["1"])
		td.Nil(synced.Items["2"])
		td.CmpDeeply(synced.Items["3"].value, 3)
		td.CmpDeeply(synced.Items["4"].value, 4)
		td.CmpDeeply(synced.Items["5"].value, 5)
		td.CmpDeeply(synced.Items["6"].value, 6)

	})

	// reverse
	c.MaxSize = 5
	c.CleanCount = 2
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(6)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(5)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(4)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(3)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(2)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(1)}
		return synced
	})

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 4)
		td.CmpDeeply(synced.Items["1"].value, 1)
		td.CmpDeeply(synced.Items["2"].value, 2)
		td.CmpDeeply(synced.Items["3"].value, 3)
		td.CmpDeeply(synced.Items["4"].value, 4)
		td.Nil(synced.Items["5"])
		td.Nil(synced.Items["6"])
	})

	c.MaxSize = 5
	c.CleanCount = 5

	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(6)}
		return synced
	})

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 0)
	})

	c.MaxSize = 5
	c.CleanCount = 6
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(6)}
		return synced
	})

	c.clean()
	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 0)
	})

	// update used time on get
	c.MaxSize = 5
	c.CleanCount = 3
	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(1)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(2)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(3)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(4)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(5)}
		synced.Items["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: newUint64Atomic(6)}
		return synced
	})

	_, _ = c.Get(ctx, "6")
	_, _ = c.Get(ctx, "2")
	_, _ = c.Get(ctx, "3")
	_, _ = c.Get(ctx, "5")
	_, _ = c.Get(ctx, "1")
	_, _ = c.Get(ctx, "4")

	c.clean()

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 3)
		td.Nil(synced.Items["6"])
		td.Nil(synced.Items["2"])
		td.Nil(synced.Items["3"])
		td.CmpDeeply(synced.Items["5"].value, 5)
		td.CmpDeeply(synced.Items["1"].value, 1)
		td.CmpDeeply(synced.Items["4"].value, 4)
	})
}

func TestLimitValueRenumberItems(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	var c = NewMemoryValueLRU("test")

	c.mu.Lock(func(synced memoryValueLRUSynced) memoryValueLRUSynced {
		synced.Items = make(map[string]*memoryValueLRUItem)
		synced.Items["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: newUint64Atomic(100)}
		synced.Items["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: newUint64Atomic(200)}
		synced.Items["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: newUint64Atomic(300)}
		synced.Items["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: newUint64Atomic(400)}
		synced.Items["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: newUint64Atomic(500)}
		return synced
	})

	c.lastTime = math.MaxUint64/2 - 1
	_ = c.Put(ctx, "6", 6)
	time.Sleep(time.Millisecond * 10)

	c.mu.RLock(func(synced memoryValueLRUSynced) {
		td.CmpDeeply(len(synced.Items), 6)

		td.CmpDeeply(synced.Items["1"].lastUsedTime.Load(), uint64(0))
		td.CmpDeeply(synced.Items["2"].lastUsedTime.Load(), uint64(1))
		td.CmpDeeply(synced.Items["3"].lastUsedTime.Load(), uint64(2))
		td.CmpDeeply(synced.Items["4"].lastUsedTime.Load(), uint64(3))
		td.CmpDeeply(synced.Items["5"].lastUsedTime.Load(), uint64(4))
		td.CmpDeeply(synced.Items["6"].lastUsedTime.Load(), uint64(5))
	})

}
