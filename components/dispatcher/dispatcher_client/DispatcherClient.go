package dispatcher_client

import (
	"../../../proto"
	"net"
)

type DistpatcherClient struct {
	proto.GoWorldConnection
}

func newDispatcherClient(conn net.Conn) *DistpatcherClient {
	return &DistpatcherClient{
		GoWorldConnection: proto.NewGoWorldConnection(conn),
	}
}
