//go:build amd64 && linux && unix

package tcp

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func control(network string, addr string, c syscall.RawConn) error {
	var opErr error
	if err := c.Control(func(fd uintptr) {
		opErr = syscall.SetsockoptInt(int(fd), unix.IPPROTO_TCP, unix.TCP_QUICKACK, 1)
	}); err != nil {
		return err
	}
	return opErr
}
