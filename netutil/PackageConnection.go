package netutil

import (
	"encoding/binary"
	"sync"
	"net"
	"github.com/TonyXMH/GoWorld/gwlog"
	"fmt"
)

const(
	MAX_PACKET_SIZE = 1*1024*1024
	SIZE_FIELD_SIZE=4
	PREPAYLOAD_SIZE=SIZE_FIELD_SIZE
	MAX_PAYLOAD_LENGTH=MAX_PACKET_SIZE-PREPAYLOAD_SIZE
)

var (
	NETWORK_ENDIAN=binary.LittleEndian
	messagePool = sync.Pool{
		New: func() interface{} {
			return &Packet{}
		},
	}
)

type PacketConnection struct {
	binconn BinaryConnection
}

func NewPacketConnection(conn net.Conn)PacketConnection  {
	return PacketConnection{binconn:NewBinaryConnection(conn)}
}

type Packet struct {
	payloadLen uint32
	bytes [MAX_PACKET_SIZE]byte
}

func (p*Packet)Payload()[]byte  {
	return p.bytes[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE+p.payloadLen]
}

func (p*Packet)Release()  {
	messagePool.Put(p)
}

func (p*Packet)AppendByte(b byte)  {
	p.bytes[PREPAYLOAD_SIZE+p.payloadLen]=b
	p.payloadLen++
}

func (p*Packet)prepareSend()  {
	NETWORK_ENDIAN.PutUint32(p.bytes[:SIZE_FIELD_SIZE],p.payloadLen)
}

func allocPacket()*Packet  {
	msg:=messagePool.Get().(*Packet)
	gwlog.Debug("ALLOC %p",msg)
	return  msg
}

func (mc*PacketConnection)NewPacket() *Packet {
	return allocPacket()
}

func (mc*PacketConnection)SendPacket(packet*Packet)error  {
	packet.prepareSend()
	err:=mc.binconn.SendAll(packet.bytes[:PREPAYLOAD_SIZE+packet.payloadLen])
	return err
}

func (mc*PacketConnection)RecvPacket()(*Packet,error)  {
	packet:= allocPacket()
	payloadLenBuf:=packet.bytes[:SIZE_FIELD_SIZE]
	err:=mc.binconn.RecvAll(payloadLenBuf)
	if err != nil{
		packet.Release()
		return nil,err
	}
	payloadLen:=NETWORK_ENDIAN.Uint32(payloadLenBuf)
	packet.payloadLen=payloadLen
	if payloadLen > MAX_PAYLOAD_LENGTH{
		packet.Release()
		return nil, fmt.Errorf("message packet too large:%v",payloadLen)
	}
	err = mc.binconn.RecvAll(packet.bytes[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE+payloadLen])
	if err != nil {
		packet.Release()
		return nil,err
	}
	return packet,nil
}
