package cache

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type DiskCache struct {
	Dir string

	mu sync.RWMutex
}

func (d *DiskCache) filepath(key string) string {
	return filepath.Join(d.Dir, diskCacheSanitizeKey(key))
}

func (d *DiskCache) Get(ctx context.Context, key string) ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	zc.L(ctx).Debug("Get from disk cache", zap.String("dir", d.Dir), zap.String("key", key))

	res, err := ioutil.ReadFile(d.filepath(key))
	if os.IsNotExist(err) {
		err = ErrCacheMiss
	}

	zc.L(ctx).Debug("Got from disk cache", zap.String("dir", d.Dir), zap.String("key", key), zap.Error(err))
	return res, err
}

func (d *DiskCache) Put(ctx context.Context, key string, data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	zc.L(ctx).Debug("Put to disk cache", zap.String("dir", d.Dir), zap.String("key", key))
	err := ioutil.WriteFile(d.filepath(key), data, 0600)
	zc.L(ctx).Debug("Put to disk cache result.", zap.String("dir", d.Dir), zap.String("key", key),
		zap.Error(err))
	return err
}

func (d *DiskCache) Delete(ctx context.Context, key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	zc.L(ctx).Debug("Delete from cache", zap.String("dir", d.Dir), zap.String("key", key))
	err := os.Remove(d.filepath(key))
	zc.L(ctx).Debug("Delete from cache result", zap.String("dir", d.Dir), zap.String("key", key),
		zap.Error(err))

	if os.IsNotExist(err) {
		err = nil
	}
	return err
}

func diskCacheSanitizeKey(k string) string {
	const placeholder = "___"
	k = strings.Replace(k, "/", placeholder, -1)
	k = strings.Replace(k, "\\", placeholder, -1)
	k = strings.Replace(k, ":", placeholder, -1)
	k = strings.Replace(k, "\"", placeholder, -1)
	k = strings.Replace(k, "'", placeholder, -1)
	k = strings.Replace(k, " ", placeholder, -1)
	return k
}
