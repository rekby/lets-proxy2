package cache

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rekby/lets-proxy2/internal/manager"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestDiskCache(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	dirPath, err := ioutil.TempDir("", "lets-proxy2-test-")
	defer os.RemoveAll(dirPath)

	if err != nil {
		t.Fatal(err)
	}

	d := &DiskCache{Dir: dirPath}
	res, err := d.Get(ctx, "asd")
	if len(res) != 0 {
		t.Error(res)
	}
	if err != manager.ErrCacheMiss {
		t.Error(err)
	}

	data := []byte("aaa")
	err = d.Put(ctx, "asd", data)
	if err != nil {
		t.Error(err)
	}

	res, err = d.Get(ctx, "asd")
	if !bytes.Equal(res, data) {
		t.Error(res)
	}
	if err != nil {
		t.Error(err)
	}

	err = d.Delete(ctx, "asd")
	if err != nil {
		t.Error(err)
	}

	res, err = d.Get(ctx, "asd")
	if len(res) != 0 {
		t.Error(res)
	}
	if err != manager.ErrCacheMiss {
		t.Error(err)
	}
}
