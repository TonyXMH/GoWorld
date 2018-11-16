package proto

import (
	. "../common"
	"../netutil"
	"net"
)

type GoWorldConnection struct {
	packetConn netutil.PacketConnection
}

func NewGoWorldConnection(conn net.Conn) GoWorldConnection {
	return GoWorldConnection{
		packetConn: netutil.NewPacketConnection(conn),
	}
}

func (gwc *GoWorldConnection) SendSetGameID(id int) error {
	packet := gwc.packetConn.NewPacket()
	packet.AppendUint16(MT_SET_GAME_ID)
	packet.AppendUint16(uint16(id))
	return gwc.packetConn.SendPacket(packet)
}
func (gwc *GoWorldConnection) SendNotifyCreateEntity(id EntityID) error {
	packet := gwc.packetConn.NewPacket()
	packet.AppendUint16(MT_NOTIFY_CREATE_ENTITY)
	packet.AppendEntityID(id)
	return gwc.packetConn.SendPacket(packet)
}

func (gwc *GoWorldConnection) SendDeclareService(id EntityID, serviceName string) error {
	packet := gwc.packetConn.NewPacket()
	packet.AppendUint16(MT_DECLEARE_SERVICE)
	packet.AppendEntityID(id)
	packet.AppendVarStr(serviceName)
	return gwc.sendPacketRelease(packet)
}

func (gwc *GoWorldConnection) SendCallEntityMethod(id EntityID, method string) error {
	packet := gwc.packetConn.NewPacket()
	packet.AppendUint16(MT_CALL_ENTITY_METHOD)
	packet.AppendEntityID(id)
	packet.AppendVarStr(method)
	return gwc.sendPacketRelease(packet)
}

func (gwc *GoWorldConnection) SendPackt(pkt *netutil.Packet) error {
	return gwc.packetConn.SendPacket(pkt)
}

func (gwc *GoWorldConnection) sendPacketRelease(pkt *netutil.Packet) error {
	err := gwc.packetConn.SendPacket(pkt)
	pkt.Release()
	return err
}
func (gwc *GoWorldConnection) Recv(msgtype *MsgType_t) (*netutil.Packet, error) {
	pkt, err := gwc.packetConn.RecvPacket()
	if err != nil {
		return nil, err
	}
	//payload := pkt.Payload()
	//*msgtype = MsgType_t(netutil.PACKET_ENDIAN.Uint16(payload[:MSG_TYPE_SIZE]))
	//*data = payload[MSG_TYPE_SIZE:]
	*msgtype = MsgType_t(pkt.ReadUint16())
	return pkt, nil
}

func (gwc *GoWorldConnection) Close() {
	gwc.packetConn.Close()
}

func (gwc *GoWorldConnection) ReomteAddr() net.Addr {
	return gwc.packetConn.RemoteAddr()
}

func (gwc *GoWorldConnection) LocalAddr() net.Addr {
	return gwc.packetConn.LocalAddr()
}

func (gwc *GoWorldConnection) String() string {
	return gwc.packetConn.String()
}
