//go:build windows

package tcp

import (
	"syscall"
)

func control(network string, addr string, c syscall.RawConn) error {
	// Windows does not support TCP_QUICKACK, so this can be a no-op.
	return nil
}
