package dispatcher

import (
	"fmt"
	"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/netutil"
	"net"
)

type DispatcherService struct {
	config *config.DispatcherConfig
}

func NewDispatcherService(cfg *config.DispatcherConfig) *DispatcherService {
	return &DispatcherService{
		config: cfg,
	}
}

func (ds *DispatcherService) Run() {
	host := fmt.Sprintf("%s:%d", ds.config.Ip, ds.config.Port)
	netutil.ServeTCPForever(host, ds)
}

func (ds *DispatcherService) ServeTCPConnection(conn net.Conn) {
	client := NewDispatcherClientProxy(conn)
	client.Serve()
}
