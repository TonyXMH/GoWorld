package proto

import (
	"bytes"
	"encoding/json"
)

type JSONMsgPacker struct {
}

func (mp JSONMsgPacker) PackMsg(msg interface{}, buf []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(buf)
	jsonEncoder := json.NewEncoder(buffer)
	err := jsonEncoder.Encode(msg)
	if err != nil {
		return buf, err
	}
	buf = buffer.Bytes()
	return buf[:len(buf)-1], nil
}

func (mp JSONMsgPacker) UnpackMsg(buf []byte, msg interface{}) error {
	return json.Unmarshal(buf, msg)
}
