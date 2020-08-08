package dockerfile // import "github.com/docker/docker/builder/dockerfile"

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	"github.com/sirupsen/logrus"
)

// ImageProber exposes an Image cache to the Builder. It supports resetting a
// cache.
//
// ImageProber向生成器公开Image缓存。
// 它支持重置缓存。
//
type ImageProber interface {
	Reset()
	Probe(parentID string, runConfig *container.Config) (string, error)
}

type imageProber struct {
	cache       builder.ImageCache
	reset       func() builder.ImageCache
	cacheBusted bool
}

func newImageProber(cacheBuilder builder.ImageCacheBuilder, cacheFrom []string, noCache bool) ImageProber {
	if noCache {
		return &nopProber{}
	}

	reset := func() builder.ImageCache {
		return cacheBuilder.MakeImageCache(cacheFrom)
	}
	return &imageProber{cache: reset(), reset: reset}
}

func (c *imageProber) Reset() {
	c.cache = c.reset()
	c.cacheBusted = false
}

// Probe checks if cache match can be found for current build instruction.
// It returns the cachedID if there is a hit, and the empty string on miss
//
// 探测检查是否可以找到当前生成指令的缓存匹配。
// 如果命中，则返回cachedID，如果未命中，则返回空字符串
//
func (c *imageProber) Probe(parentID string, runConfig *container.Config) (string, error) {
	if c.cacheBusted {
		return "", nil
	}
	cacheID, err := c.cache.GetCache(parentID, runConfig)
	if err != nil {
		return "", err
	}
	if len(cacheID) == 0 {
		// docker build 第八个log
		logrus.Debugf("[BUILDER] Cache miss: %s", runConfig.Cmd)
		c.cacheBusted = true
		return "", nil
	}
	logrus.Debugf("[BUILDER] Use cached version: %s", runConfig.Cmd)
	return cacheID, nil
}

type nopProber struct{}

func (c *nopProber) Reset() {}

func (c *nopProber) Probe(_ string, _ *container.Config) (string, error) {
	return "", nil
}
