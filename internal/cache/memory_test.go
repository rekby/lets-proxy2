package cache

import (
	"bytes"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestMemoryCache(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	c := NewMemoryCache("test")
	res, err := c.Get(ctx, "asd")
	if len(res) != 0 {
		t.Error(res)
	}
	if err != ErrCacheMiss {
		t.Error(err)
	}

	data := []byte("aaa")
	err = c.Put(ctx, "asd", data)
	e.CmpNoError(err)

	res, err = c.Get(ctx, "asd")
	if !bytes.Equal(res, data) {
		t.Error(res)
	}
	e.CmpNoError(err)

	err = c.Delete(ctx, "asd")
	e.CmpNoError(err)

	err = c.Delete(ctx, "non-existed-key")
	e.CmpNoError(err)

	res, err = c.Get(ctx, "asd")
	if len(res) != 0 {
		t.Error(res)
	}
	if err != ErrCacheMiss {
		t.Error(err)
	}
}
