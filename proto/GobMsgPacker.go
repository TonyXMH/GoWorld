package proto

import (
	"bytes"
	"encoding/gob"
)

//encoding/gob是go自带的数据结构编解码的工具

type GobMsgPacker struct {
}

func (mp GobMsgPacker) PackMsg(msg interface{}, buf []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(buf)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(msg)
	if err != nil {
		return buf, err
	}
	buf = buffer.Bytes()
	return buf, err
}

func (mp GobMsgPacker) UnpackMsg(data []byte, msg interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	err := decoder.Decode(msg)
	return err
}
