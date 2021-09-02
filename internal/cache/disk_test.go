package cache

import (
	"bytes"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestDiskCache(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	dirPath := th.TmpDir(e)

	c := &DiskCache{Dir: dirPath}
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

	err = c.Delete(ctx, "non-existed-key")
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
