package images // import "github.com/docker/docker/daemon/images"

import (
	"github.com/docker/docker/builder"
	"github.com/docker/docker/image/cache"
	"github.com/sirupsen/logrus"
)

// MakeImageCache creates a stateful image cache.
// MakeImageCache创建状态Image缓存。
func (i *ImageService) MakeImageCache(sourceRefs []string) builder.ImageCache {
	if len(sourceRefs) == 0 {
		return cache.NewLocal(i.imageStore)
	}

	cache := cache.New(i.imageStore)

	for _, ref := range sourceRefs {
		img, err := i.GetImage(ref, nil)
		if err != nil {
			logrus.Warnf("Could not look up %s for cache resolution, skipping: %+v", ref, err)
			continue
		}
		cache.Populate(img)
	}

	return cache
}
