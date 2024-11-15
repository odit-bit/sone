package tcp

import (
	"context"
	"net"
)

func Listen(ctx context.Context, addr string) (net.Listener, error) {
	lc := net.ListenConfig{Control: control}
	return lc.Listen(ctx, "tcp", addr)

}
