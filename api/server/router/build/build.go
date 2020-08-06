package build // import "github.com/docker/docker/api/server/router/build"

import (
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/api/types"
)

// buildRouter is a router to talk with the build controller
type buildRouter struct {
	backend  Backend
	daemon   experimentalProvider
	routes   []router.Route
	features *map[string]bool
}

// NewRouter initializes a new build router
func NewRouter(b Backend, d experimentalProvider, features *map[string]bool) router.Router {
	r := &buildRouter{
		backend:  b,
		daemon:   d,
		features: features,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the build controller
func (r *buildRouter) Routes() []router.Route {
	return r.routes
}

func (r *buildRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewPostRoute("/build", r.postBuild),
		router.NewPostRoute("/build/prune", r.postPrune),
		router.NewPostRoute("/build/cancel", r.postCancel),
	}
}

// BuilderVersion derives the default docker builder version from the config
// Note: it is valid to have BuilderVersion unset which means it is up to the
// client to choose which builder to use.
//
// BuilderVersion从配置派生默认的停靠器构建器版本注意：取消设置BuilderVersion是有效的，这意味着由客户端选择使用哪个构建器。
//
func BuilderVersion(features map[string]bool) types.BuilderVersion {
	var bv types.BuilderVersion
	if v, ok := features["buildkit"]; ok {
		if v {
			bv = types.BuilderBuildKit
		} else {
			bv = types.BuilderV1
		}
	}
	return bv
}
