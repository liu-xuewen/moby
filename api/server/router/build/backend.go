package build // import "github.com/docker/docker/api/server/router/build"

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
)

// Backend abstracts an image builder whose only purpose is to build an image referenced by an imageID.
// 后端抽象出一个镜像构建器，该构建器的唯一目的是构建由imageID引用的镜像。
type Backend interface {
	// Build a Docker image returning the id of the image
	// TODO: make this return a reference instead of string
	Build(context.Context, backend.BuildConfig) (string, error)

	// Prune build cache
	PruneCache(context.Context, types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error)

	Cancel(context.Context, string) error
}

type experimentalProvider interface {
	HasExperimental() bool
}
