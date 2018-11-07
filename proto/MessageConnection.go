package proto

import (
	"../gwlog"
	"../netutil"
	"../uuid"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

const (
	MAX_MESSAGE_SIZE      = 1 * 1024 * 1024
	SIZE_FIELD_SIZE       = 4
	TYPE_FIELD_SIZE       = 2
	PREPAYLOAD_SIZE       = SIZE_FIELD_SIZE + TYPE_FIELD_SIZE
	STRING_ID_SIZE        = uuid.UUID_LENGTH
	RELAY_PREPAYLOAD_SIZE = SIZE_FIELD_SIZE + STRING_ID_SIZE + TYPE_FIELD_SIZE
	RELAY_MASK            = 0x80000000
)

var (
	NETWORK_ENDIAN = binary.LittleEndian
	messagePool    = sync.Pool{
		New: newMessageInPool,
	}
)

func newMessageInPool() interface{} {
	return &Message{}
}

type MessageConnection struct {
	netutil.BinaryConnection
}

func NewMessageConnection(conn net.Conn) MessageConnection {
	return MessageConnection{netutil.NewBinaryConnection(conn)}
}

type Message [MAX_MESSAGE_SIZE]byte

func allocMessage() *Message {
	msg := messagePool.Get().(*Message)
	gwlog.Debug("ALLOC %p", msg)
	return msg
}

func (m *Message) Release() {
	messagePool.Put(m)
}

func toJsonString(msg interface{}) string {
	s, _ := json.Marshal(messagePool)
	return string(s)
}

//[size 4B][type 2B][payload NB]
func (mc *MessageConnection) SendMsg(mt MsgType_t, msg interface{}) error {
	return mc.SendMsgEx(mt, msg, MSG_PACKER)
}

//打包发送
func (mc *MessageConnection) SendMsgEx(mt MsgType_t, msg interface{}, msgPacker MsgPacker) error {
	msgbuf := allocMessage()
	defer msgbuf.Release()

	NETWORK_ENDIAN.PutUint16(msgbuf[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := msgbuf[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE] //6:end
	payloadCap := cap(payloadBuf)
	payloadBuf, err := msgPacker.PackMsg(msg, payloadBuf)
	if err != nil {
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		return fmt.Errorf("MessageConnection: message payload too large(%d):%v", payloadLen, msg)
	}
	pktSize := uint32(payloadLen + PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32(msgbuf[:SIZE_FIELD_SIZE], pktSize)
	err = mc.SendAll((msgbuf)[:pktSize])
	gwlog.Debug(">>> SendMsg: size=%v, %s%v,error=%v", pktSize, MsgTypeToString(mt), toJsonString(msg), err)
	return err

}

//[size4B][stringID][type 2B][payload NB]
func (mc *MessageConnection) SendRelayMsg(targetID string, mt MsgType_t, msg interface{}) error {
	msgbuf := allocMessage()
	defer msgbuf.Release()
	copy(msgbuf[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+STRING_ID_SIZE], []byte(targetID))
	NETWORK_ENDIAN.PutUint16(msgbuf[SIZE_FIELD_SIZE+STRING_ID_SIZE:SIZE_FIELD_SIZE+STRING_ID_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := msgbuf[RELAY_PREPAYLOAD_SIZE:RELAY_PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := MSG_PACKER.PackMsg(msg, payloadBuf)
	if err != nil {
		return err
	}
	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		return fmt.Errorf("MessageConnection: message payload too large(%d):%v", payloadLen, msg)
	}
	pktSize := uint32(payloadLen + RELAY_PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32(msgbuf[:SIZE_FIELD_SIZE], pktSize|RELAY_MASK)
	err = mc.SendAll((msgbuf)[:pktSize])
	gwlog.Debug(">>> SendRelayMsg: size=%v, %s%v,error=%v", pktSize, MsgTypeToString(mt), toJsonString(msg), err)
	return err
}

type MessageHandler interface {
	HandleMsg(msg *Message, pktSize uint32, msgType MsgType_t) error
	HandleRelayMsg(msg *Message, pktSize uint32, targetID string) error
}

func (mc *MessageConnection) RecvMsg(handler MessageHandler) error {
	msg := allocMessage()
	pktSizeBuf := msg[:SIZE_FIELD_SIZE]
	err := mc.RecvAll(pktSizeBuf)
	if err != nil {
		return err
	}
	pktSize := NETWORK_ENDIAN.Uint32(pktSizeBuf)
	isRelayMsg := false
	if pktSize&RELAY_MASK != 0 {
		isRelayMsg = true
		pktSize -= RELAY_MASK
	}
	if pktSize > MAX_MESSAGE_SIZE {
		msg.Release()
		return fmt.Errorf("message packet too large:%v", pktSize)
	}
	err = mc.RecvAll(msg[SIZE_FIELD_SIZE:pktSize])
	if err != nil {
		msg.Release()
		return err
	}
	gwlog.Debug("<<<RecvMsg:pktsize=%v,isRelayMsg=%v,packet=%v", pktSize, isRelayMsg, msg[:pktSize])
	if isRelayMsg {
		targetID := string(msg[SIZE_FIELD_SIZE : SIZE_FIELD_SIZE+STRING_ID_SIZE])
		err = handler.HandleRelayMsg(msg, pktSize, targetID)
	} else {
		msgtype := MsgType_t(NETWORK_ENDIAN.Uint16(msg[SIZE_FIELD_SIZE : SIZE_FIELD_SIZE+TYPE_FIELD_SIZE]))
		err = handler.HandleMsg(msg, pktSize, msgtype)
	}
	return err
}
