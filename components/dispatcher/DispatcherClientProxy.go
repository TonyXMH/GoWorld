package dispatcher

import "net"

type DispatcherClientProxy struct {
	net.Conn
}

func NewDispatcherClientProxy(conn net.Conn) *DispatcherClientProxy {
	return &DispatcherClientProxy{conn}
}

func (dcp *DispatcherClientProxy) Serve() {
	dcp.Close()
}
