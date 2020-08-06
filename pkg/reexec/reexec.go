package reexec // import "github.com/docker/docker/pkg/reexec"

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var registeredInitializers = make(map[string]func())

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	}

	registeredInitializers[name] = initializer
}

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.
// init作为exec进程的第一部分调用，如果调用了初始化函数，则返回true。
func Init() bool {
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()

		return true
	}
	return false
}

func naiveSelf() string {
	name := os.Args[0]
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err == nil {
			return lp
		}
	}
	// handle conversion of relative paths to absolute
	if absName, err := filepath.Abs(name); err == nil {
		return absName
	}
	// if we couldn't get absolute name, return original
	// (NOTE: Go only errors on Abs() if os.Getwd fails)
	return name
}
