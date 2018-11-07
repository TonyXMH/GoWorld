package dispatcher

import (
	"../../entity"
	"../../netutil"
	"../../proto"
	"fmt"
	"github.com/TonyXMH/GoWorld/gwlog"
	"net"
)

type DispatcherClientProxy struct {
	proto.GoWorldConnection
	owner  *DispatcherService
	gameid int
}

func NewDispatcherClientProxy(owner *DispatcherService, conn net.Conn) *DispatcherClientProxy {
	return &DispatcherClientProxy{
		GoWorldConnection: proto.NewGoWorldConnection(conn),
		owner:             owner,
	}
}

func (dcp *DispatcherClientProxy) Serve() {
	defer func() {
		dcp.Close()
		err := recover()
		if err != nil && !netutil.IsConnectionClosed(err) {
			gwlog.Error("Client %s paniced with error: %v", dcp, err)
		}
	}()
	gwlog.Info("New dispatcher client: %s", dcp)
	for {
		var msgtype proto.MsgType_t
		pkt, err := dcp.Recv(&msgtype)
		if err != nil {
			gwlog.Panic(err)
		}
		gwlog.Info("%s.RecvPacket: msgtype=%v,payload=%v", dcp, msgtype, pkt.Payload())
		if msgtype == proto.MT_NOTIFY_CREATE_ENTITY {
			eid := entity.EntityID(pkt.ReadBytes(entity.ENTITY_LENGTH))
			dcp.owner.HandleNotifyCreateEntity(dcp, eid)
		} else if msgtype == proto.MT_SET_GAME_ID {
			gameid := int(pkt.ReadUint16())
			dcp.gameid = gameid
			dcp.owner.HandleSetGameID(dcp, gameid)
		}
	}
}

func (dcp *DispatcherClientProxy) String() string {
	return fmt.Sprintf("DispatcherClientProxy<%d|>", dcp.gameid)
	//return fmt.Sprintf("DispatcherClientProxy<%d|%s>", dcp.gameid, dcp.RemoteAddr())
}
