package cache

import (
	"context"
	"github.com/rekby/lets-proxy2/internal/log"
	"github.com/rekby/safemutex"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
)

type DiskCache struct {
	Dir string

	mu safemutex.RWMutex[struct{}]
}

func (c *DiskCache) filepath(key string) string {
	return filepath.Join(c.Dir, diskCacheSanitizeKey(key))
}

func (c *DiskCache) Get(ctx context.Context, key string) (res []byte, err error) {
	c.mu.RLock(func(synced struct{}) {
		filePath := c.filepath(key)

		logLevel := zapcore.DebugLevel
		res, err = os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				err = ErrCacheMiss
			} else {
				logLevel = zapcore.ErrorLevel
			}
		}

		log.LevelParamCtx(ctx, logLevel, "Got from disk cache", zap.String("dir", c.Dir), zap.String("key", key),
			zap.String("file", filePath), zap.Error(err))
	})

	return res, err
}

func (c *DiskCache) Put(ctx context.Context, key string, data []byte) (err error) {
	c.mu.Lock(func(synced struct{}) struct{} {
		zc.L(ctx).Debug("Put to disk cache", zap.String("dir", c.Dir), zap.String("key", key))
		err = os.WriteFile(c.filepath(key), data, 0600)
		zc.L(ctx).Debug("Put to disk cache result.", zap.String("dir", c.Dir), zap.String("key", key),
			zap.Error(err))

		return synced
	})

	return err
}

func (c *DiskCache) Delete(ctx context.Context, key string) (err error) {
	c.mu.Lock(func(synced struct{}) struct{} {
		zc.L(ctx).Debug("Delete from cache", zap.String("dir", c.Dir), zap.String("key", key))
		err = os.Remove(c.filepath(key))
		zc.L(ctx).Debug("Delete from cache result", zap.String("dir", c.Dir), zap.String("key", key),
			zap.Error(err))

		if os.IsNotExist(err) {
			err = nil
		}

		return synced
	})

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
