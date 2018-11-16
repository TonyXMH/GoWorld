package dispatcher

import (
	"../../common"
	"../../config"
	"../../gwlog"
	"../../netutil"
	"fmt"
	"net"
)

type DispatcherService struct {
	config             *config.DispatcherConfig
	clients            []*DispatcherClientProxy
	entityLocs         map[common.EntityID]int
	registeredServices map[string]common.EntityID
}

func NewDispatcherService(cfg *config.DispatcherConfig) *DispatcherService {
	return &DispatcherService{
		config:             cfg,
		clients:            []*DispatcherClientProxy{},
		entityLocs:         map[common.EntityID]int{},
		registeredServices: map[string]common.EntityID{},
	}
}

func (service *DispatcherService) String() string {
	return fmt.Sprintf("DispatchersService<C%d|E%d>", len(service.clients), len(service.entityLocs))
}

func (service *DispatcherService) Run() {
	host := fmt.Sprintf("%s:%d", service.config.Ip, service.config.Port)
	netutil.ServeTCPForever(host, service)
}

func (service *DispatcherService) ServeTCPConnection(conn net.Conn) {
	client := NewDispatcherClientProxy(service, conn)
	client.Serve()
}

func (service *DispatcherService) HandleSetGameID(dcp *DispatcherClientProxy, pkt *netutil.Packet, gameid int) {
	gwlog.Debug("%s.HandleSetGameID:dcp=%s, gameid=%d", service, dcp, gameid)
	for gameid >= len(service.clients) {
		service.clients = append(service.clients, nil)
	}
	service.clients[gameid] = dcp
	pkt.Release()
	return
}

func (service *DispatcherService) HandleNotifyCreateEntity(dcp *DispatcherClientProxy, pkt *netutil.Packet, entityID common.EntityID) {
	gwlog.Debug("%s.HandleNotifyCreateEntity:dcp=%s,entity=%s", service, dcp, entityID)
	service.entityLocs[entityID] = dcp.gameid
	pkt.Release()
}

func (service *DispatcherService) HandleDeclareService(dcp *DispatcherClientProxy, pkt netutil.Packet, entityID common.EntityID, serviceName string) {
	gwlog.Debug("%s.HandleDeclareService:dcp=%s,entityID=%s,serviceName=%s", service, dcp, entityID, serviceName)
	service.broadcastToDispatcherClients(pkt)
	pkt.Release()
}

func (service *DispatcherService) broadcastToDispatcherClients(pkt *netutil.Packet) {
	for _, dcp := range service.clients {
		if dcp != nil {
			dcp.SendPacket(pkt)
		}
	}
}
