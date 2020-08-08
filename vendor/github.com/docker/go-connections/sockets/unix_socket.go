// +build !windows

package sockets

import (
	"net"
	"os"
	"syscall"
)

// NewUnixSocket creates a unix socket with the specified path and group.
func NewUnixSocket(path string, gid int) (net.Listener, error) {
	// 从文件系统中删除一个名称。如果名称是文件的最后一个连接，并且没有其它进程将文件打开，名称对应的文件会实际被删除。
	// IsNotExist返回一个布尔值，指示是否已知错误以报告文件或目录不存在。
	// ErrNotExist以及一些系统调用错误都能满足它。
	if err := syscall.Unlink(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	mask := syscall.Umask(0777)
	defer syscall.Umask(mask)

	// 直接listen unix文件会自动创建unix监听文件，该文件存在则报错
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chown(path, 0, gid); err != nil {
		l.Close()
		return nil, err
	}
	if err := os.Chmod(path, 0660); err != nil {
		l.Close()
		return nil, err
	}
	return l, nil
}
