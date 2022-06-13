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
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 6}
	c.clean()
	td.CmpDeeply(len(c.m), 6)
	td.CmpDeeply(c.m["1"].value, 1)
	td.CmpDeeply(c.m["2"].value, 2)
	td.CmpDeeply(c.m["3"].value, 3)
	td.CmpDeeply(c.m["4"].value, 4)
	td.CmpDeeply(c.m["5"].value, 5)
	td.CmpDeeply(c.m["6"].value, 6)

	c.MaxSize = 5
	c.CleanCount = 3
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.clean()
	td.CmpDeeply(len(c.m), 5)
	td.CmpDeeply(c.m["1"].value, 1)
	td.CmpDeeply(c.m["2"].value, 2)
	td.CmpDeeply(c.m["3"].value, 3)
	td.CmpDeeply(c.m["4"].value, 4)
	td.CmpDeeply(c.m["5"].value, 5)

	c.MaxSize = 5
	c.CleanCount = 2
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 6}
	c.clean()
	td.CmpDeeply(len(c.m), 4)
	td.Nil(c.m["1"])
	td.Nil(c.m["2"])
	td.CmpDeeply(c.m["3"].value, 3)
	td.CmpDeeply(c.m["4"].value, 4)
	td.CmpDeeply(c.m["5"].value, 5)
	td.CmpDeeply(c.m["6"].value, 6)

	// reverse
	c.MaxSize = 5
	c.CleanCount = 2
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 6}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 5}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 4}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 3}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 2}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 1}
	c.clean()
	td.CmpDeeply(len(c.m), 4)
	td.CmpDeeply(c.m["1"].value, 1)
	td.CmpDeeply(c.m["2"].value, 2)
	td.CmpDeeply(c.m["3"].value, 3)
	td.CmpDeeply(c.m["4"].value, 4)
	td.Nil(c.m["5"])
	td.Nil(c.m["6"])

	c.MaxSize = 5
	c.CleanCount = 5
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 6}
	c.clean()
	td.CmpDeeply(len(c.m), 0)

	c.MaxSize = 5
	c.CleanCount = 6
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 6}
	c.clean()
	td.CmpDeeply(len(c.m), 0)

	// update used time on get
	c.MaxSize = 5
	c.CleanCount = 3
	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 1}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 2}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 3}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 4}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 5}
	c.m["6"] = &memoryValueLRUItem{key: "6", value: 6, lastUsedTime: 6}
	_, _ = c.Get(ctx, "6")
	_, _ = c.Get(ctx, "2")
	_, _ = c.Get(ctx, "3")
	_, _ = c.Get(ctx, "5")
	_, _ = c.Get(ctx, "1")
	_, _ = c.Get(ctx, "4")
	c.clean()
	td.CmpDeeply(len(c.m), 3)
	td.Nil(c.m["6"])
	td.Nil(c.m["2"])
	td.Nil(c.m["3"])
	td.CmpDeeply(c.m["5"].value, 5)
	td.CmpDeeply(c.m["1"].value, 1)
	td.CmpDeeply(c.m["4"].value, 4)
}

func TestLimitValueRenumberItems(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	var c = NewMemoryValueLRU("test")

	c.m = make(map[string]*memoryValueLRUItem)
	c.m["1"] = &memoryValueLRUItem{key: "1", value: 1, lastUsedTime: 100}
	c.m["2"] = &memoryValueLRUItem{key: "2", value: 2, lastUsedTime: 200}
	c.m["3"] = &memoryValueLRUItem{key: "3", value: 3, lastUsedTime: 300}
	c.m["4"] = &memoryValueLRUItem{key: "4", value: 4, lastUsedTime: 400}
	c.m["5"] = &memoryValueLRUItem{key: "5", value: 5, lastUsedTime: 500}

	c.lastTime = math.MaxUint64/2 - 1
	_ = c.Put(ctx, "6", 6)
	time.Sleep(time.Millisecond * 10)
	td.CmpDeeply(len(c.m), 6)

	c.mu.RLock()
	defer c.mu.RLock()
	td.CmpDeeply(c.m["1"].lastUsedTime, uint64(0))
	td.CmpDeeply(c.m["2"].lastUsedTime, uint64(1))
	td.CmpDeeply(c.m["3"].lastUsedTime, uint64(2))
	td.CmpDeeply(c.m["4"].lastUsedTime, uint64(3))
	td.CmpDeeply(c.m["5"].lastUsedTime, uint64(4))
	td.CmpDeeply(c.m["6"].lastUsedTime, uint64(5))
}
