package dispatcher

import (
	"../../entity"
	"fmt"
	"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/gwlog"
	"github.com/TonyXMH/GoWorld/netutil"
	"net"
)

type DispatcherService struct {
	config     *config.DispatcherConfig
	clients    map[int]*DispatcherClientProxy
	entityLocs map[entity.EntityID]int
}

func NewDispatcherService(cfg *config.DispatcherConfig) *DispatcherService {
	return &DispatcherService{
		config:     cfg,
		entityLocs: map[entity.EntityID]int{},
	}
}

func (ds *DispatcherService) String() string {
	return fmt.Sprintf("DispatchersService<C%d|E%d>", len(ds.clients), len(ds.entityLocs))
}

func (ds *DispatcherService) Run() {
	host := fmt.Sprintf("%s:%d", ds.config.Ip, ds.config.Port)
	netutil.ServeTCPForever(host, ds)
}

func (ds *DispatcherService) ServeTCPConnection(conn net.Conn) {
	client := NewDispatcherClientProxy(ds, conn)
	client.Serve()
}

func (ds *DispatcherService) HandleSetGameID(dcp *DispatcherClientProxy, gameid int) {
	gwlog.Debug("%s.HandleSetGameID:dcp=%s, gameid=%d", ds, dcp, gameid)
	return
}

func (ds *DispatcherService) HandleNotifyCreateEntity(dcp *DispatcherClientProxy, entityID entity.EntityID) {
	gwlog.Debug("%s.HandleNotifyCreateEntity:dcp=%s,entity=%s", ds, dcp, entityID)
	ds.entityLocs[entityID] = dcp.gameid
}
