package gate

import (
	"fmt"
	"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/netutil"
	"net"
)

type GateServer struct {
	config *config.GateConfig
}

func NewGateServer(gatecfg *config.GateConfig) *GateServer {
	return &GateServer{config: gatecfg}
}

func (gs *GateServer) Run() {
	listenAddr := fmt.Sprintf("%s:%d", gs.config.Ip, gs.config.Port)
	netutil.ServeTCPForever(listenAddr, gs)
}

func (gs *GateServer) ServeTCPConnection(conn net.Conn) {
	clientProxy := newGateClientProxy(conn)
	clientProxy.serve()
}
