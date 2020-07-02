package cache

import (
	"context"
	"github.com/rekby/lets-proxy2/internal/log"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type DiskCache struct {
	Dir string

	mu sync.RWMutex
}

func (c *DiskCache) filepath(key string) string {
	return filepath.Join(c.Dir, diskCacheSanitizeKey(key))
}

func (c *DiskCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := c.filepath(key)

	logLevel := zapcore.DebugLevel
	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrCacheMiss
		} else {
			logLevel = zapcore.ErrorLevel
		}
	}

	log.LevelParamCtx(ctx, logLevel, "Got from disk cache", zap.String("dir", c.Dir), zap.String("key", key),
		zap.String("file", filePath), zap.Error(err))

	return res, err
}

func (c *DiskCache) Put(ctx context.Context, key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	zc.L(ctx).Debug("Put to disk cache", zap.String("dir", c.Dir), zap.String("key", key))
	err := ioutil.WriteFile(c.filepath(key), data, 0600)
	zc.L(ctx).Debug("Put to disk cache result.", zap.String("dir", c.Dir), zap.String("key", key),
		zap.Error(err))
	return err
}

func (c *DiskCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	zc.L(ctx).Debug("Delete from cache", zap.String("dir", c.Dir), zap.String("key", key))
	err := os.Remove(c.filepath(key))
	zc.L(ctx).Debug("Delete from cache result", zap.String("dir", c.Dir), zap.String("key", key),
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
