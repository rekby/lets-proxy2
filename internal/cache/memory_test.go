package cache

import (
	"bytes"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestMemoryCache(t *testing.T) {
	ctx, flush := th.TestContext()
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
	if err != nil {
		t.Error(err)
	}

	res, err = c.Get(ctx, "asd")
	if !bytes.Equal(res, data) {
		t.Error(res)
	}
	if err != nil {
		t.Error(err)
	}

	err = c.Delete(ctx, "asd")
	if err != nil {
		t.Error(err)
	}

	res, err = c.Get(ctx, "asd")
	if len(res) != 0 {
		t.Error(res)
	}
	if err != ErrCacheMiss {
		t.Error(err)
	}
}
