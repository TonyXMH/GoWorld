package uuid

import (
	"encoding/base64"
	"encoding/binary"
	"time"
	"os"
	"io"
	"crypto/rand"
	"fmt"
	"crypto/md5"
	"sync/atomic"
)

const (
	UUID_LENGTH = 16
	encodeUUID  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_."
)
//
var UUIDEncoding=base64.NewEncoding(encodeUUID).WithPadding(base64.NoPadding)
var objectIdCnt uint32=0
var machineIds=readMachineIds()
func GenUUID() string {
	b:=make([]byte,12)
	binary.BigEndian.PutUint32(b[:],uint32(time.Now().Unix()))
	b[4]=machineIds[0]
	b[5]=machineIds[1]
	b[6]=machineIds[2]
	pid:=os.Getpid()
	b[7]=byte(pid>>8)
	b[8]=byte(pid)
	i:=atomic.AddUint32(&objectIdCnt,1)
	b[9]=byte(i>>16)
	b[10]=byte(i>>8)
	b[11]=byte(i)
	return UUIDEncoding.EncodeToString(b)
}


func readMachineIds() []byte {
	ids:=make([]byte,3)
	hostname,err1:=os.Hostname()
	if err1!=nil{
		_,err2:=io.ReadFull(rand.Reader,ids)
		if err2!=nil{
			panic(fmt.Errorf("cannot get hostname:%v;%v",err1,err2))
		}
		return ids
	}
	hw:=md5.New()
	hw.Write([]byte(hostname))
	copy(ids,hw.Sum(nil))
	return ids
}
